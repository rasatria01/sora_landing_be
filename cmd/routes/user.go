package routes

import (
	"sora_landing_be/cmd/constants"
	"sora_landing_be/cmd/controllers"
	"sora_landing_be/cmd/services"
	"sora_landing_be/pkg/http/server/middlewares"

	"github.com/gin-gonic/gin"
)

func registerUser(router *gin.RouterGroup) {
	userCtl := controllers.NewUserController(services.ServicePool.UserService)

	user := router.Group("/users")
	{
		user.GET("", middlewares.RoleHandler(constants.UserRoleSuperAdmin), userCtl.ListUser)
		user.GET("profile", userCtl.GetProfile)
		user.GET(":id", middlewares.RoleHandler(constants.UserRoleSuperAdmin), userCtl.Get)
		user.POST("", middlewares.RoleHandler(constants.UserRoleSuperAdmin), userCtl.CreateUser)
		user.PUT(":id", middlewares.RoleHandler(constants.UserRoleSuperAdmin), userCtl.Update)
		user.DELETE(":id", middlewares.RoleHandler(constants.UserRoleSuperAdmin), userCtl.Delete)
	}
}
