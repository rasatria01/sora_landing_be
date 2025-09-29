package services

import (
	"context"
	"database/sql"
	"errors"
	"sora_landing_be/cmd/constants"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/cmd/dto/response"
	"sora_landing_be/cmd/repository"
	"sora_landing_be/pkg/authentication"
	"sora_landing_be/pkg/database"
	internal_err "sora_landing_be/pkg/errors"

	"github.com/uptrace/bun"
)

type AuthService interface {
	Login(ctx context.Context, payload requests.Login) (res response.LoginResponse, err error)
	Logout(ctx context.Context) (err error)
	RefreshToken(ctx context.Context, payload requests.RefreshToken) (res response.RefreshTokenResponse, err error)
}

type authService struct {
	authRepo repository.AuthRepository
}

func NewAuthSrv(authRepo repository.AuthRepository) AuthService {
	return &authService{authRepo: authRepo}
}

func (a *authService) Login(ctx context.Context, payload requests.Login) (res response.LoginResponse, err error) {
	err = database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		authData, err := a.authRepo.GetByEmail(ctx, payload.Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return internal_err.AuthError(constants.AuthPasswordInvalidOrEmailNotFound)
			}
			return err
		}

		isValid, err := authentication.VerifyPassword(payload.Password, authData.Password)
		if err != nil {
			return err
		}

		if !isValid {
			return internal_err.AuthError(constants.AuthPasswordInvalidOrEmailNotFound)
		}

		tokenPayload := requests.ToTokenPayload(authData)
		pair, err := authentication.JWTAuth.GenerateTokenPair(tokenPayload, false)
		if err != nil {
			return err
		}

		var authUpdated = authData
		authUpdated.RefreshTokenID = &tokenPayload.RefreshTokenID
		err = a.authRepo.UpdateAuth(ctx, &authUpdated)
		if err != nil {
			return err
		}

		res = response.LoginResponse{
			AccessToken:  pair.AccessToken,
			RefreshToken: pair.RefreshToken,
		}

		return nil
	})
	if err != nil {
		return res, err
	}

	return res, nil
}

func (a *authService) Logout(ctx context.Context) (err error) {
	err = database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		authToken := authentication.GetUserDataFromToken(ctx)
		if authToken.AuthID == "" {
			return internal_err.AuthError(internal_err.DataNotFound)
		}

		authData, err := a.authRepo.GetByID(ctx, &authToken.AuthID, nil)
		if err != nil {
			return err
		}

		authData.RefreshTokenID = nil
		err = a.authRepo.UpdateAuth(ctx, &authData)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *authService) RefreshToken(ctx context.Context, payload requests.RefreshToken) (res response.RefreshTokenResponse, err error) {
	claims, err := authentication.JWTAuth.VerifyRefreshToken(payload.RefreshToken)
	if err != nil {
		return res, err
	}
	if claims == nil {
		return res, internal_err.AuthError(constants.AuthInvalidToken)
	}

	auth, err := a.authRepo.GetByID(ctx, nil, &claims.TokenID)
	if err != nil { 
		return res, err
	}
	tokenPayload := requests.ToTokenPayload(auth)
	claimsRefresh, err := authentication.JWTAuth.GenerateTokenPair(tokenPayload, true)
	if err != nil {
		return res, err
	}

	return response.RefreshTokenResponse{
		AccessToken: claimsRefresh.AccessToken,
	}, err
}
