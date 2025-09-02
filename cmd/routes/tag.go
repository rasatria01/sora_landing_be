package routes

import (
	"sora_landing_be/cmd/controllers"
	"sora_landing_be/cmd/services"

	"github.com/gin-gonic/gin"
)

func registerTag(router *gin.RouterGroup) {
	userCtl := controllers.NewTagController(services.ServicePool.TagService)

	tag := router.Group("/tag")
	{
		tag.POST("", userCtl.Create)
		tag.GET("", userCtl.List)
		tag.GET(":id", userCtl.Get)
		tag.PUT(":id", userCtl.Update)
		tag.DELETE(":id", userCtl.Delete)
	}
}
