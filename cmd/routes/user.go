package routes

import (
	"sora_landing_be/cmd/controllers"
	"sora_landing_be/cmd/services"

	"github.com/gin-gonic/gin"
)

func registerUser(router *gin.RouterGroup) {
	userCtl := controllers.NewUserController(services.ServicePool.UserService)

	user := router.Group("/users")
	{
		user.GET("", userCtl.ListUser)
		user.GET("profile", userCtl.GetProfile)
		user.GET(":id", userCtl.Get)
		user.PUT(":id", userCtl.Update)
		user.DELETE(":id", userCtl.Delete)
	}
}
