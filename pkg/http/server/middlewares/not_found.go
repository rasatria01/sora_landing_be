package middlewares

import (
	"net/http"
	"sora_landing_be/pkg/errors"
	"sora_landing_be/pkg/http/server/http_response"

	"github.com/gin-gonic/gin"
)

func NotFoundHandler(c *gin.Context) {
	http_response.SendError(c, errors.NewDefaultError(http.StatusNotFound, "Not Found"))
}
