package handlers

import (
	"clean_api/src/api/dto"
	"clean_api/src/api/helpers"
	"clean_api/src/api/validators"
	"clean_api/src/services"
	"errors"
	"net/http"
	"time"
)

type UsersHandler struct {
	service *services.UserService
	tk *services.TokenService
}

func NewUsersHandler() *UsersHandler{
	u := services.NewUserService()
	t := services.NewTokenService()
	return &UsersHandler{service: u, tk: t}
}

func (uh *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	req := dto.CreateUser{}
	err := helpers.ReadRequestBody(w, r, &req)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}
	
	v := validators.NewValidator()

	if validators.ValidateUser(v, &req); !v.Valid(){
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(v.Errors, false, http.StatusBadRequest,errors.New("validation error")))
		return
	}

	res, err := uh.service.Create(&req)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}

	tk, err := uh.tk.New(res.ID, time.Hour*6, services.ScopeActivation)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}


	helpers.WriteResponse(w, helpers.GenerateResponse(map[string]interface{}{"user":res, "token":tk}, false, http.StatusBadRequest))


}

func (uh *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	req := dto.UpdateUser{}
	err := helpers.ReadRequestBody(w, r, &req)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}

	id, err := helpers.ReadParams(r)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}

	user, err := uh.service.GetById(int(id))
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}

	v := validators.NewValidator()
	if req.Email != "" {
		if validators.ValidateEmail(v, req.Email);!v.Valid(){
			helpers.WriteResponse(w, helpers.GenerateResponseWithError(v.Errors, false, http.StatusBadRequest,errors.New("validation error")))
			return
		}
	} else {
		req.Email = user.Email
	}
	if req.Name == ""{
		req.Name = user.Name
	}
	if !req.Activated {
		req.Activated = user.Activated
	}

	res, err := uh.service.Update(&req, int(id), user.Version)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}
	helpers.WriteResponse(w, helpers.GenerateResponse(res, false, http.StatusBadRequest))

}

func (uh *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadParams(r)
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}
	res, err := uh.service.GetById(int(id))
	if err != nil {
		helpers.WriteResponse(w, helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err))
		return
	}
	helpers.WriteResponse(w, helpers.GenerateResponse(res, false, http.StatusBadRequest))

}
