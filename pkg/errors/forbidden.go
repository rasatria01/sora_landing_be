package errors

import (
	"net/http"
)

func NewForbiddenError(message string, err string) AppError {
	return AppError{
		Code:    http.StatusForbidden,
		Message: message,
		Err:     err,
	}
}

func ForbiddenErrorToAppError() AppError {
	return NewForbiddenError("Forbidden access", "User has not access to this resource")
}
