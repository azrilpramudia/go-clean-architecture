package model

// request user register
type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	Name string `json:"name" validate:"required,max=100"`
}

// response to client
type UserRespone struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
	Name string `json:"name"`
}

type LoginUserRequest struct {
	Username string `json:"Username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type TokenResponse struct {
	Token string `json:"token"`
}