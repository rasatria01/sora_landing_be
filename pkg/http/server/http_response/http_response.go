package http_response

import (
	"errors"
	"net/http"
	internal_err "sora_landing_be/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   any    `json:"error"`
}

func SendSuccess(c *gin.Context, status int, message string, data any) {
	var msg string
	if message != "" {
		msg = message
	} else {
		msg = http.StatusText(status)
	}

	response := Response{
		Status:  status,
		Message: msg,
		Data:    data,
	}
	c.JSON(status, response)
}

func SendError(c *gin.Context, err error) {
	var appErr internal_err.AppError
	if errors.As(err, &appErr) {
		err = c.Error(appErr)
	} else {
		err = c.Error(err)
	}

	if err != nil {
		c.Abort()
	}
}

func SendRaw(c *gin.Context, status int, contentType string, body string) {
	c.Header("Content-Type", contentType)
	c.String(status, body)
}
