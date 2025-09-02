package controllers

import (
	"net/http"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/cmd/services"
	"sora_landing_be/pkg/errors"
	internalHTTP "sora_landing_be/pkg/http"
	"sora_landing_be/pkg/http/server/http_response"

	"github.com/gin-gonic/gin"
)

type TagController struct {
	TagServices services.TagService
}

func NewTagController(tagServices services.TagService) TagController {
	return TagController{
		TagServices: tagServices,
	}
}

func (ctl *TagController) Create(ctx *gin.Context) {
	var payload requests.TagRequest
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err := ctl.TagServices.CreateTag(ctx, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusCreated, "", nil)
}

func (ctl *TagController) List(ctx *gin.Context) {
	var branchs requests.ListTag
	if err := internalHTTP.BindData(ctx, &branchs); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.TagServices.ListTag(ctx, branchs)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}
	http_response.SendSuccess(ctx, http.StatusOK, "Success get list tag", res)
}
func (ctl *TagController) Update(ctx *gin.Context) {
	var payload requests.TagRequest
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err = internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.TagServices.UpdateTag(ctx, id, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}
func (ctl *TagController) Delete(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.TagServices.DeleteTag(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}

func (ctl *TagController) Get(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.TagServices.GetTag(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get data", res)
}
