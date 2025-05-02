package errors

import "errors"

var (
	ErrValidateToken = errors.New("unexpected signing method")
	ErrUserNotFound  = errors.New("user not found")
)
