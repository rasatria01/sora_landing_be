package errors

import (
	"errors"
	"net/http"
	"sora_landing_be/pkg/utils"

	"github.com/go-playground/validator/v10"
)

func NewValidationError(message string, err map[string]string) AppError {
	return AppError{
		Code:    http.StatusBadRequest, // HTTP status code for bad request
		Message: message,
		Err:     err, // Validation details
	}
}

func ValidationErrorToAppError(err error) AppError {
	validationErrors := make(map[string]string)

	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		for _, e := range errs {
			validationErrors[utils.ToSnakeCase(e.Field())] = validationErrorToText(e)
		}
	} else {
		return NewDefaultError(http.StatusBadRequest, err.Error())
	}

	return NewValidationError("Validation failed", validationErrors)
}

var errorWording = map[string]string{
	"required":         "is required",
	"email":            "is an invalid email address",
	"is_digit":         "is not a digit",
	"is_only_alphabet": "is not a valid alphabet",
	"min":              "is too short",
	"max":              "is too long",
}

func validationErrorToText(e validator.FieldError) string {
	if wording, ok := errorWording[e.Tag()]; ok {
		return wording
	}
	return "invalid"
}
