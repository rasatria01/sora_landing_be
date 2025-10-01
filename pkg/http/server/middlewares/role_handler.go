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
		userRoles := authentication.GetUserRoleFromToken(c) // assume returns map[UserRole]bool

		// check if any of the allowed roles is in user's roles
		for _, role := range roles {
			if userRoles[role] {
				// authorized → continue
				c.Next()
				return
			}
		}

		// if no match, reject
		http_response.SendError(c, errors.ForbiddenErrorToAppError())
		c.Abort()
	}
}
