package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func isUppercase(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == value[:]
}

func isPureString(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	re := regexp.MustCompile(CheckOnlyAlphabet)
	return re.MatchString(value)
}

func isDigit(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	re := regexp.MustCompile(Digits)
	return re.MatchString(value)
}


func validateEnum(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(Enum); ok {
		return value.IsValidEnum()
	}
	return false
}
