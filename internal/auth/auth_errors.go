package auth

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidEmail   = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrEmailExists    = errors.New("email already exists")
)
