package routes

import (
	"sora_landing_be/cmd/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterSeo(router *gin.Engine) {
	seoCtl := controllers.NewSeoController()

	seo := router.Group("/")
	{
		seo.GET("robots.txt", seoCtl.GetRobots)
	}
}
