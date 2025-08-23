package controllers

import (
	"net/http"
	"sora_landing_be/cmd/constants"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/cmd/services"
	"sora_landing_be/pkg/errors"
	internalHTTP "sora_landing_be/pkg/http"
	"sora_landing_be/pkg/http/server/http_response"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authSrv services.AuthService) AuthController {
	return AuthController{
		authService: authSrv,
	}
}

func (ctl *AuthController) Login(ctx *gin.Context) {
	var auth requests.Login
	if err := internalHTTP.BindData(ctx, &auth); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.authService.Login(ctx, auth)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, constants.AuthLoginSuccess, res)
}

func (ctl *AuthController) Logout(ctx *gin.Context) {
	err := ctl.authService.Logout(ctx)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, constants.AuthLogoutSuccess, nil)
}

func (ctl *AuthController) RefreshToken(ctx *gin.Context) {
	var auth requests.RefreshToken
	if err := internalHTTP.BindData(ctx, &auth); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.authService.RefreshToken(ctx, auth)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, constants.AuthRefreshTokenSuccess, res)
}
