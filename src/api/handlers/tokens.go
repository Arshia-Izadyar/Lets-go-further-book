package handlers

import (
	"clean_api/src/api/dto"
	"clean_api/src/api/helpers"
	"clean_api/src/api/validators"
	"clean_api/src/services"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TokenHandler struct {
	TkService *services.TokenService
	UsrService *services.UserService
}

func NewTokenHandler() *TokenHandler{
	s := services.NewTokenService()
	u := services.NewUserService()
	return &TokenHandler{
		TkService:  s,
		UsrService: u,
	}
}

func (th *TokenHandler) CreateAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	req := &dto.AuthTokenRequest{}
	err := helpers.ReadRequestBody(w, r, &req)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}

	v := validators.NewValidator()
	validators.ValidateEmail(v, req.Email)
	validators.ValidatePasswordPlaintext(v, req.Password)
	if !v.Valid(){
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(v.Errors, false, http.StatusBadRequest, nil))
		return
	}
	user, err := th.UsrService.GetByEmail(req.Email)

	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}
	token, err := th.TkService.New(user.ID, time.Hour*6, services.ScopeAuthentication)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}
	helpers.WriteResponse(w, helpers.GenerateResponse(token, true, 201))
}


func (th *TokenHandler) ActiveUser(w http.ResponseWriter, r *http.Request) {
	req := &dto.ActiveUserDTO{}
	err := helpers.ReadRequestBody(w, r, &req)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}

	v := validators.NewValidator()
	if validators.ValidateTokenPlainText(v, req.Token); !v.Valid(){
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(v.Errors, false, http.StatusBadRequest,errors.New("validation error")))
		return
	}

	u, err := th.TkService.GetForToken(services.ScopeActivation, req.Token)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}

	updateDTO := &dto.UpdateUser{
		Name:      u.Name,
		Email:     u.Email,
		Activated: true,
	}
	fmt.Println(u)
	usr, err := th.UsrService.Update(updateDTO, int(u.ID), u.Version)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}
	err = th.TkService.DeleteAllForUser(services.ScopeActivation, int(u.ID))
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}
	helpers.WriteResponse(w, helpers.GenerateResponse(usr, true, 201))

}

func (th *TokenHandler) ResetAccessToken(w http.ResponseWriter, r *http.Request) {
	req := &dto.AuthTokenRequest{}
	err := helpers.ReadRequestBody(w, r, &req)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}

	v := validators.NewValidator()
	validators.ValidateEmail(v, req.Email)
	validators.ValidatePasswordPlaintext(v, req.Password)
	if !v.Valid(){
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(v.Errors, false, http.StatusBadRequest, nil))
		return
	}
	user, err := th.UsrService.GetByEmail(req.Email)

	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}
	err = th.TkService.DeleteAllForUser(services.ScopeAuthentication, user.ID)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}
	
	helpers.WriteResponse(w, helpers.GenerateResponse("deleted", true, 201))
}