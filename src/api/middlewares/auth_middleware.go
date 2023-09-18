package middlewares

import (
	"clean_api/src/api/helpers"
	"clean_api/src/api/validators"
	"clean_api/src/data/models"
	"clean_api/src/services"
	"errors"
	"net/http"
	"strings"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Authorization")
		authorizeHeader := r.Header.Get("Authorization")
		if authorizeHeader == "" {
			r := helpers.ContextSetUser(r, models.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authorizeHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer"{
			res := helpers.GenerateResponseWithError(nil, false, http.StatusUnauthorized, errors.New("invalid authentication"))
			helpers.WriteResponse(w, res)
			return
		}
		token := headerParts[1]
		v := validators.NewValidator()
		if validators.ValidateTokenPlainText(v, token); !v.Valid(){
			res := helpers.GenerateResponseWithError(nil, false, http.StatusUnauthorized, errors.New("invalid authentication"))
			helpers.WriteResponse(w, res)
			return
		}
		tkService := services.NewTokenService()
		user, err := tkService.GetForToken(services.ScopeAuthentication, token)
		if err != nil {
			res := helpers.GenerateResponseWithError(nil, false, http.StatusUnauthorized, err)
			helpers.WriteResponse(w, res)
			return
		}
		r = helpers.ContextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func RequireActiveUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := helpers.ContextGetUser(r)
		isAnon := helpers.IsAnon(&user)
		if isAnon {
			res := helpers.GenerateResponseWithError(nil, false, http.StatusUnauthorized, errors.New("authentication required"))
			helpers.WriteResponse(w, res)
			return
		}
		if !user.Activated {
			res := helpers.GenerateResponseWithError(nil, false, http.StatusUnauthorized, errors.New("user activation required"))
			helpers.WriteResponse(w, res)
			return
		}
		next.ServeHTTP(w, r)
	})
}


func RequirePermission(code string,next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := helpers.ContextGetUser(r)

		permService := services.NewPermissionService()

		perms, err := permService.GetAllForUser(int(user.ID))
		if err != nil {
			res := helpers.GenerateResponseWithError(nil, false, http.StatusUnauthorized, err)
			helpers.WriteResponse(w, res)
			return
		}
		for _, p := range perms {
			if p == code{
				next.ServeHTTP(w, r)
				return
			}
		}
		res := helpers.GenerateResponseWithError(nil, false, http.StatusUnauthorized, errors.New("user unauthorized user dont have required role"))
		helpers.WriteResponse(w, res)
		return
			
	})
}