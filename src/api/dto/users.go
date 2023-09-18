package dto

import "time"

type CreateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string	`json:"password"`
}

type UpdateUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Activated bool `json:"activated"`
}


type GetUser struct {
	Email string `json:"email"`
}

type UserResponse struct {
	ID        int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Activated bool `json:"activated"`
	Password string `json:"-"`
	Version   int `json:"version"`
}