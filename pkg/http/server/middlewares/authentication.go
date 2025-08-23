package middlewares

import (
	"context"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/pkg/authentication"
	internal_err "sora_landing_be/pkg/errors"
	"sora_landing_be/pkg/http/server/http_response"
	"sora_landing_be/pkg/logger"

	"strings"

	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			http_response.SendError(c, internal_err.AuthError(authentication.AuthErrMalformedToken.Error()))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http_response.SendError(c, internal_err.AuthError(authentication.AuthErrInvalidToken.Error()))
			return
		}

		accessToken, err := authentication.JWTAuth.VerifyAccessToken(tokenString)
		if err != nil {
			http_response.SendError(c, err)
			return
		}
		userAuth := requests.UserAuth{
			AuthID:         accessToken.AuthID,
			UserID:         accessToken.UserID,
			Email:          accessToken.Email,
			Role:           accessToken.Role,
			RefreshTokenID: accessToken.RefreshTokenID,
		}

		c.Set(authentication.Token, userAuth)
		logger.Log.Debug(accessToken.UserID)
		ctx := context.WithValue(c.Request.Context(), authentication.Token, userAuth)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
