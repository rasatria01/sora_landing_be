package routes

import (
	"sora_landing_be/cmd/controllers"
	"sora_landing_be/cmd/services"

	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(router *gin.RouterGroup) {
	userCtl := controllers.NewFileController(services.ServicePool.FileService)

	publicFiles := router.Group("/files")
	{
		publicFiles.POST("", userCtl.UploadFiles)
		publicFiles.GET("/:id", userCtl.GetPublicFile)
		publicFiles.DELETE("/:id", userCtl.DeleteFile)
	}
}
