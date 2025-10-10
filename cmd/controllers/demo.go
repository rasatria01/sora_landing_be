package controllers

import (
	"fmt"
	"net/http"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/cmd/services"
	"sora_landing_be/pkg/errors"
	internalHTTP "sora_landing_be/pkg/http"
	"sora_landing_be/pkg/http/server/http_response"

	"github.com/gin-gonic/gin"
)

type DemoController struct {
	demoServices services.DemoService
}

func NewDemoController(demoService services.DemoService) DemoController {
	return DemoController{
		demoServices: demoService,
	}
}

func (ctl *DemoController) Create(ctx *gin.Context) {
	var payload requests.Demo
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err := ctl.demoServices.CreateDemo(ctx, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusCreated, "", nil)
}

func (ctl *DemoController) List(ctx *gin.Context) {
	var branchs requests.ListDemo
	if err := internalHTTP.BindData(ctx, &branchs); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.demoServices.ListDemos(ctx, branchs)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}
	http_response.SendSuccess(ctx, http.StatusOK, "Success get list demo", res)
}
func (ctl *DemoController) Update(ctx *gin.Context) {
	var payload requests.Demo
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err = internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.demoServices.UpdateDemo(ctx, id, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}
func (ctl *DemoController) Delete(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.demoServices.DeleteDemo(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}

func (ctl *DemoController) Get(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.demoServices.GetDemoByID(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get data", res)
}

func (ctl *DemoController) Export(ctx *gin.Context) {
	var payload requests.ExportDemo
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
	}

	file, err := ctl.demoServices.ExportDemo(ctx, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}
	buf, err := file.WriteToBuffer()
	if err != nil {
		http_response.SendError(ctx, err)

		return
	}
	defer file.Close()

	ctx.Header("Content-Disposition", "attachment; filename=export-demo.xlsx")
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))
	if err := file.Write(ctx.Writer); err != nil {
		http_response.SendError(ctx, errors.StorageErrorToAppError("failed to write file"))
		return
	}
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
