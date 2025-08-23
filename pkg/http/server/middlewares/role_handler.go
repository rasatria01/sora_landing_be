package middlewares

import (
	"sora_landing_be/cmd/constants"
	"sora_landing_be/pkg/authentication"
	"sora_landing_be/pkg/errors"
	"sora_landing_be/pkg/http/server/http_response"

	"github.com/gin-gonic/gin"
)

func RoleHandler(roles ...constants.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, role := range roles {
			if authentication.GetUserRoleFromToken(c)[role] {
				break
			} else {
				http_response.SendError(c, errors.ForbiddenErrorToAppError())
			}
		}

		c.Next()
	}
}
