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

type UserController struct {
	UserService services.UserService
}

func NewUserController(bankSrv services.UserService) UserController {
	return UserController{
		UserService: bankSrv,
	}
}

func (ctl *UserController) CreateUser(ctx *gin.Context) {
	var agent requests.CreateUser
	if err := internalHTTP.BindData(ctx, &agent); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err := ctl.UserService.Register(ctx, agent)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusCreated, "", nil)
}

func (ctl *UserController) ListUser(ctx *gin.Context) {
	var users requests.ListUser
	if err := internalHTTP.BindData(ctx, &users); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.UserService.GetList(ctx, users)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get list data", res)
}

func (ctl *UserController) Update(ctx *gin.Context) {
	var vendor requests.CreateUser

	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err = internalHTTP.BindData(ctx, &vendor); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.UserService.Update(ctx, id, vendor)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}

func (ctl *UserController) Delete(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.UserService.DeleteSrv(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}

func (ctl *UserController) Get(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.UserService.Detail(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get data", res)
}

func (ctl *UserController) GetProfile(ctx *gin.Context) {
	res, err := ctl.UserService.Detail(ctx, authentication.GetUserDataFromToken(ctx).UserID)
	if err != nil {
		http_response.SendError(ctx, err)

		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get data", res)
}
