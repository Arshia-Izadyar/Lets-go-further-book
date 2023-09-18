package helpers

import (
	"clean_api/src/data/models"
	"context"
	"net/http"
)

type ContestKey string

const UserContextKey = ContestKey("user")

func ContextGetUser(r *http.Request) models.Users {
	users, ok := r.Context().Value(UserContextKey).(*models.Users)
	if !ok {
		panic("missing user value in ctx")
	}
	return *users
}

func ContextSetUser(r *http.Request, user *models.Users) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func IsAnon(u *models.Users) bool{
	return u == models.AnonymousUser
}