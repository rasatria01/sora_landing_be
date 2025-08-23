package client

type HTTPClientError struct {
	StatusCode int
	Message    string
}

func (e *HTTPClientError) Error() string {
	return e.Message
}

func NewHTTPClientError(statusCode int, message string) *HTTPClientError {
	return &HTTPClientError{
		StatusCode: statusCode,
		Message:    message,
	}
}
