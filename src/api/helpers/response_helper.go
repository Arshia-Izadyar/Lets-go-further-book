package helpers

import (
	"clean_api/src/api/validators"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Response struct {
	Success    bool        `json:"success"`
	Result     interface{} `json:"result"`
	StatusCode int         `json:"statusCode"`
	Err        string      `json:"err"`
}

func GenerateResponse(result interface{}, success bool, code int) *Response {
	return &Response{
		Success:    success,
		Result:     result,
		StatusCode: code,
	}
}

func GenerateResponseWithError(result interface{}, success bool, code int, err error) *Response {
	return &Response{
		Success:    success,
		Result:     result,
		StatusCode: code,
		Err:        err.Error(),
	}
}

func WriteResponse(w http.ResponseWriter, response *Response) error {
	jsonData, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return err
	}
	w.WriteHeader(response.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func ReadRequestBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	max := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(max))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("syntax error at body %d", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("badly formed json")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect type %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect at %q", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body is empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("must send only one json")
	}
	return nil
}


func ReadParams(r *http.Request) (int64, error){
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func ReadString(qs *url.Values, key, defaultValue string) string{
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

func ReadInt(qs *url.Values, key string, defaultValue int, v *validators.Validator) int{
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an int")
		return defaultValue
	}
	return i
}


func ReadCSV(qs *url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, ",")
}