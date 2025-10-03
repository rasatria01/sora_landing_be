package routes

import (
	"sora_landing_be/cmd/controllers"
	"sora_landing_be/cmd/services"

	"github.com/gin-gonic/gin"
)

func registerPublic(router *gin.RouterGroup) {
	ctl := controllers.NewDemoController(services.ServicePool.DemoService)
	bctl := controllers.NewBlogController(services.ServicePool.BlogService)

	payment := router.Group("/demo")
	{
		payment.POST("/create", ctl.Create)
	}

	blog := router.Group("/blog")
	{
		blog.GET(":id", bctl.GetPublicArticleBySlug)
		blog.GET("", bctl.ListPublicArticles)
		blog.GET("/featured", bctl.GetFeaturedArticle)
	}

}
