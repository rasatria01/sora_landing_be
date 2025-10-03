package routes

import (
	"sora_landing_be/cmd/controllers"
	"sora_landing_be/cmd/services"

	"github.com/gin-gonic/gin"
)

func registerBlog(router *gin.RouterGroup) {
	blogCtl := controllers.NewBlogController(services.ServicePool.BlogService)

	blog := router.Group("/articles")
	{
		// Read operations
		blog.GET("", blogCtl.ListArticles)
		blog.GET("stats", blogCtl.GetArticleStats)
		blog.GET(":id", blogCtl.GetArticle)
		blog.GET("by-slug/:slug", blogCtl.GetArticleBySlug)

		// Write operations
		blog.POST("", blogCtl.CreateArticle)
		blog.POST("external", blogCtl.CreateArticleFromURL)
		blog.PUT(":id", blogCtl.UpdateArticle)
		blog.PATCH(":id/status", blogCtl.UpdateArticleStatus)
		blog.PUT(":id/tags", blogCtl.UpdateArticleTags)
		blog.PATCH(":id/set-featured", blogCtl.SetFeaturedPosition)
		blog.PATCH(":id/del-featured", blogCtl.RemoveFeaturedPosition)

		// Delete operations
		blog.DELETE(":id", blogCtl.DeleteArticle)
		blog.DELETE(":id/hard", blogCtl.HardDeleteArticle)
	}
}
