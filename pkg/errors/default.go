package errors

type AppError struct {
	Code    int         `json:"-"`
	Message string      `json:"-"`
	Err     interface{} `json:"-"`
}

func (e AppError) Error() string {
	return e.Message
}

func NewDefaultError(code int, message string) AppError {
	return AppError{
		Code:    code,
		Message: message,
	}
}
