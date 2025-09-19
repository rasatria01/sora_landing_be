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

type CatController struct {
	CatServices services.CategoryService
}

func NewCatController(catServices services.CategoryService) CatController {
	return CatController{
		CatServices: catServices,
	}
}

func (ctl *CatController) Create(ctx *gin.Context) {
	var payload requests.Category
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err := ctl.CatServices.CreateCategory(ctx, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusCreated, "", nil)
}

func (ctl *CatController) List(ctx *gin.Context) {
	var branchs requests.ListCategory
	if err := internalHTTP.BindData(ctx, &branchs); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.CatServices.ListCategory(ctx, branchs)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}
	http_response.SendSuccess(ctx, http.StatusOK, "Success get list tag", res)
}
func (ctl *CatController) Update(ctx *gin.Context) {
	var payload requests.Category
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err = internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.CatServices.UpdateCategory(ctx, id, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}
func (ctl *CatController) Delete(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.CatServices.DeleteCategory(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}

func (ctl *CatController) Get(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.CatServices.GetCategory(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get data", res)
}
