package routes

import (
	"net/http"
	"sora_landing_be/cmd/repository"
	"sora_landing_be/cmd/services"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/http/server/http_response"
	"sora_landing_be/pkg/http/server/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterV1(router *gin.Engine) {
	RegisterSeo(router)

	v1 := router.Group("/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {

			http_response.SendSuccess(c, http.StatusOK, "pong", gin.H{
				"time": time.Now(),
				"ua":   c.Request.UserAgent(),
			})
		})

		repository.Init(database.GetDB())
		services.Init()

		registerAuth(v1)
		registerTag(v1)
		v1.Use(middlewares.TokenAuthMiddleware())
		registerCategory(v1)
		registerUser(v1)
		registerBlog(v1)

	}

}
