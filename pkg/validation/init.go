package validation

import (
	"log"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type validatorFunc struct {
	validation    validator.Func
	callIfDataNil bool
}

var validatorList = make(map[string]validatorFunc)
var once sync.Once

func InitGinValidator() {
	once.Do(func() {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			initCustomValidator()
			for key, validationFunc := range validatorList {
				if err := v.RegisterValidation(key, validationFunc.validation, validationFunc.callIfDataNil); err != nil {
					log.Fatalf("Error registering %s validator: %v", key, err)
				}
			}
		} else {
			log.Fatal("failed to retrieve validator engine")
		}
	})
}

func registerValidator(key string, fn validator.Func, callValidationIfNil bool) {
	validatorList[key] = validatorFunc{
		validation:    fn,
		callIfDataNil: callValidationIfNil,
	}
}

func initCustomValidator() {
	registerValidator("is_upper_case", isUppercase, false)
	registerValidator("is_digit", isDigit, false)
	registerValidator("valid_enum", validateEnum, false)
	registerValidator("is_only_alphabet", isPureString, false)
}
