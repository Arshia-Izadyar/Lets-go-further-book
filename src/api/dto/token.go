package dto

type AuthTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ActiveUserDTO struct {
	Token string `json:"token"`
}