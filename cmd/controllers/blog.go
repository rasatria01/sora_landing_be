package controllers

import (
	"net/http"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/cmd/services"
	"sora_landing_be/pkg/authentication"
	"sora_landing_be/pkg/errors"
	internalHTTP "sora_landing_be/pkg/http"
	"sora_landing_be/pkg/http/server/http_response"

	"github.com/gin-gonic/gin"
)

type BlogController struct {
	BlogService services.BlogService
}

func NewBlogController(blogService services.BlogService) BlogController {
	return BlogController{
		BlogService: blogService,
	}
}

func (ctl *BlogController) CreateArticle(ctx *gin.Context) {
	var payload requests.BlogArtikel
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	userID := authentication.GetUserDataFromToken(ctx).UserID
	err := ctl.BlogService.CreateArticle(ctx, userID, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusCreated, "Article created successfully", nil)
}

func (ctl *BlogController) CreateArticleFromURL(ctx *gin.Context) {
	var payload requests.FromURL
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	userID := authentication.GetUserDataFromToken(ctx).UserID
	err := ctl.BlogService.CreateArticleFromURL(ctx, userID, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusCreated, "Article created successfully", nil)
}

func (ctl *BlogController) UpdateArticle(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	var payload requests.UpdateArtikel
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.BlogService.UpdateArticle(ctx, id, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Article updated successfully", nil)
}

func (ctl *BlogController) UpdateArticleStatus(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	var payload requests.UpdateArticleStatus
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.BlogService.UpdateArticleStatus(ctx, id, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Article status updated successfully", nil)
}

func (ctl *BlogController) GetArticle(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	article, err := ctl.BlogService.GetArticle(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Article retrieved successfully", article)
}

func (ctl *BlogController) GetArticleBySlug(ctx *gin.Context) {
	slug, err := internalHTTP.BindParams[string](ctx, "slug")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	article, err := ctl.BlogService.GetArticleBySlug(ctx, slug)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Article retrieved successfully", article)
}

func (ctl *BlogController) ListArticles(ctx *gin.Context) {
	var params requests.ListArtikel
	if err := internalHTTP.BindData(ctx, &params); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	articles, err := ctl.BlogService.ListArticles(ctx, params)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Articles retrieved successfully", articles)
}

func (ctl *BlogController) GetArticleStats(ctx *gin.Context) {
	stats, err := ctl.BlogService.GetArticleStats(ctx)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Article statistics retrieved successfully", stats)
}

func (ctl *BlogController) UpdateArticleTags(ctx *gin.Context) {
	articleID, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	var payload struct {
		TagIDs []string `json:"tag_ids" binding:"required"`
	}
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.BlogService.UpdateArticleTags(ctx, articleID, payload.TagIDs)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Article tags updated successfully", nil)
}

func (ctl *BlogController) DeleteArticle(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.BlogService.DeleteArticle(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Article deleted successfully", nil)
}

func (ctl *BlogController) HardDeleteArticle(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.BlogService.HardDeleteArticle(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Article permanently deleted successfully", nil)
}

// Public endpoints

func (ctl *BlogController) ListPublicArticles(ctx *gin.Context) {
	var params requests.ListArtikel
	if err := internalHTTP.BindData(ctx, &params); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	articles, err := ctl.BlogService.ListPublicArticles(ctx, params)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Articles retrieved successfully", articles)
}

func (ctl *BlogController) GetPublicArticleBySlug(ctx *gin.Context) {
	slug, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	article, err := ctl.BlogService.GetPublicArticleBySlug(ctx, slug)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Article retrieved successfully", article)
}
