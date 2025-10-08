package routes

import (
	"sora_landing_be/cmd/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(router *gin.RouterGroup) {
	userCtl := controllers.NewFileController()

	publicFiles := router.Group("/files")
	{
		publicFiles.POST("", userCtl.UploadFile)
		publicFiles.DELETE(":filename", userCtl.DeleteFile)
	}
}
