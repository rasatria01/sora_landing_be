package middlewares

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sora_landing_be/pkg/config"
	internal_errs "sora_landing_be/pkg/errors"
	"sora_landing_be/pkg/http/server/http_response"
	"sora_landing_be/pkg/utils"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				http_response.SendError(c, internal_errs.NewDefaultError(http.StatusInternalServerError, fmt.Sprintf("%v", r)))
				c.JSON(http.StatusInternalServerError, http_response.Response{
					Status:  http.StatusInternalServerError,
					Message: http.StatusText(http.StatusInternalServerError),
				})
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var customErr internal_errs.AppError

			switch true {
			case errors.As(err, &customErr):
				c.JSON(customErr.Code, http_response.Response{
					Status:  customErr.Code,
					Message: customErr.Message,
					Error:   customErr.Err,
				})
			case errors.Is(err, sql.ErrNoRows):
				c.JSON(http.StatusNotFound, http_response.Response{
					Status:  http.StatusNotFound,
					Message: internal_errs.DataNotFound,
				})
			default:
				c.JSON(http.StatusInternalServerError, http_response.Response{
					Status:  http.StatusInternalServerError,
					Message: http.StatusText(http.StatusInternalServerError),
					Error:   utils.Fallback(err, nil, config.LoadConfig().Application.Environment == config.Development),
				})
			}
			c.Abort()
		}
	}
}
