package services

import (
	"context"
	"database/sql"
	"sora_landing_be/cmd/constants"
	"sora_landing_be/cmd/dto"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/cmd/dto/response"
	"sora_landing_be/cmd/repository"
	"sora_landing_be/pkg/authentication"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/utils"

	"github.com/uptrace/bun"
)

type UserService interface {
	Register(ctx context.Context, payload requests.CreateUser) error
	GetList(ctx context.Context, payload requests.ListUser) (dto.PaginationResponse[response.User], error)
	Update(ctx context.Context, id string, payload requests.CreateUser) error
	DeleteSrv(ctx context.Context, userID string) error
	Detail(ctx context.Context, userID string) (response.User, error)
	Profile(ctx context.Context, userID string) (response.Profile, error)
}

type userService struct {
	userRepo repository.UserRepository
	authRepo repository.AuthRepository
}

func NewUserSrv(userRepo repository.UserRepository, authRepo repository.AuthRepository) UserService {
	return &userService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (u *userService) Register(ctx context.Context, payload requests.CreateUser) error {
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		user := payload.ToDomain()

		err := u.userRepo.CreateUser(ctx, &user)
		if err != nil {
			return err
		}

		hashedPassword, err := authentication.HashPassword(payload.Password)
		if err != nil {
			return err
		}

		auth := requests.CreateAuth{UserID: user.ID, Password: hashedPassword}.ToDomain()
		err = u.authRepo.CreateAuth(ctx, &auth)
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

func (u *userService) GetList(ctx context.Context, payload requests.ListUser) (dto.PaginationResponse[response.User], error) {
	var paginateRes dto.PaginationResponse[response.User]
	res, count, err := u.userRepo.ListUser(ctx, payload)
	if err != nil {
		return paginateRes, err
	}

	paginateRes = dto.NewPaginationResponse(payload.PaginationRequest, count, response.NewListUser(res))
	return paginateRes, nil
}

func (u *userService) Update(ctx context.Context, id string, payload requests.CreateUser) error {
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		user := payload.ToDomain()
		user.ID = id

		err := u.userRepo.UpdateUser(ctx, &user)
		if err != nil {
			return err
		}

		if payload.Password != "" {
			hashedPassword, err := authentication.HashPassword(payload.Password)
			if err != nil {
				return err
			}

			auth := requests.CreateAuth{UserID: user.ID, Password: hashedPassword}.ToDomain()
			err = u.authRepo.UpdateAuth(ctx, &auth)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) DeleteSrv(ctx context.Context, userID string) error {
	return u.userRepo.DeleteUser(ctx, userID)
}

func (u *userService) Detail(ctx context.Context, userID string) (response.User, error) {
	var res response.User
	data, err := u.userRepo.GetUser(ctx, userID)
	if err != nil {
		return res, err
	}

	res = response.NewUser(data)
	return res, nil
}
func (u *userService) Profile(ctx context.Context, userID string) (response.Profile, error) {
	var res response.Profile
	data, err := u.userRepo.GetUser(ctx, userID)
	if err != nil {
		return res, err
	}
	permit := []string{}
	if utils.Contains(data.Roles, constants.UserRoleSuperAdmin) {
		permit = []string{"users"}
	}
	res = response.NewProfile(data, permit)
	return res, nil
}
