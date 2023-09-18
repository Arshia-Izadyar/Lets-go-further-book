package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	quotedJson := strconv.Quote(fmt.Sprintf("%d mins", r))
	return []byte(quotedJson), nil
}

func (r *Runtime) UnmarshalJSON(j []byte) error {
	unquoted, err := strconv.Unquote(string(j))
	if err != nil {
		return err
	}
	parts := strings.Split(unquoted, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return errors.New("invalid format for runtime '<time> mins'")
	}
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return err
	}
	*r = Runtime(i)
	return nil
}
