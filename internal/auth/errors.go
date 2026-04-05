package auth

import "errors"

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrTokenNotValidYet   = errors.New("token not valid yet")
)
