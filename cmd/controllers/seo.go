package controllers

import (
	"net/http"
	"sora_landing_be/pkg/http/server/http_response"

	"github.com/gin-gonic/gin"
)

type SeoController struct {
}

func NewSeoController() SeoController {
	return SeoController{}
}

func (ctl *SeoController) GetRobots(ctx *gin.Context) {
	robots := `User-agent: *
Disallow: /admin
Disallow: /api
Disallow: /auth
Disallow: /static/
Allow: /

Sitemap: https://yourdomain.com/sitemap.xml`
	http_response.SendRaw(ctx, http.StatusOK, "text/plain", robots)
}
