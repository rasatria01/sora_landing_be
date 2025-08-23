package errors

import (
	"net/http"
)

func NewAuthError(message, errsMeg string) AppError {
	return AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
		Err:     errsMeg,
	}
}

func AuthError(errsMeg string) AppError {
	return NewAuthError("Authentication failed", errsMeg)
}
