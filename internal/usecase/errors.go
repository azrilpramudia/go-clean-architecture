package usecase

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUsernameAlreadyExists = errors.New("username already registered")
	ErrInvalidCredentials = errors.New("username or password is wrong")
)