package errors

import (
	"net/http"
)

func NewStorageError(message string, err string) AppError {
	return AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Err:     err,
	}
}

func StorageErrorToAppError(msg string) AppError {
	return NewStorageError("Storage validation failed", msg)
}
