package routers

import (
	"clean_api/src/api/handlers"
	"clean_api/src/api/middlewares"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Router() http.Handler {
	router := httprouter.New()
	h := handlers.NewMovieHandler()
	userHandler := handlers.NewUsersHandler()
	tk := handlers.NewTokenHandler()
	router.HandlerFunc(http.MethodPost, "/v1/movie", middlewares.RequireActiveUser(h.CreateMovie))
	router.HandlerFunc(http.MethodPut, "/v1/movie/:id", middlewares.RequireActiveUser(h.UpdateMovie))
	router.HandlerFunc(http.MethodDelete, "/v1/movie/:id", middlewares.RequireActiveUser(h.DeleteMovie))
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", middlewares.RequirePermission("movies:write",h.GetById))
	router.HandlerFunc(http.MethodGet, "/v1/movie", middlewares.RequireActiveUser(h.GetAll))


	router.HandlerFunc(http.MethodPost, "/v1/users", userHandler.CreateUser)
	router.HandlerFunc(http.MethodPut, "/v1/users/:id", userHandler.UpdateUser)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", userHandler.GetUser)

	router.HandlerFunc(http.MethodPost, "/v1/users/auth", tk.CreateAuthenticationToken)
	router.HandlerFunc(http.MethodPost, "/v1/users/active", tk.ActiveUser)
	router.HandlerFunc(http.MethodPost, "/v1/users/reset-token", tk.ResetAccessToken)


	return middlewares.RateLimiter(middlewares.Authentication(router))
}
