package middlewares

import (
	"sora_landing_be/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start).Milliseconds()
		status := c.Writer.Status()
		clientIP := c.ClientIP()

		logFields := []zap.Field{
			zap.String("path", path),
			zap.String("method", method),
			zap.Int("status", status),
			zap.Int64("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 && c.Writer.Status() >= 500 {
			logger.Log.Error("Request error", append(logFields, zap.String("errors", c.Errors.String()))...)
		} else {
			logger.Log.Info("Request completed", logFields...)
		}
	}
}
