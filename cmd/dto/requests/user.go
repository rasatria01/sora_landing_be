package requests

import (
	"sora_landing_be/cmd/constants"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/dto"
)

type (
	CreateUser struct {
		Name     string               `json:"name" binding:"required"`
		Email    string               `json:"email" binding:"required"`
		Role     []constants.UserRole `json:"role" binding:"required,dive,valid_enum"`
		Password string               `json:"password" binding:"required"`
	}

	Login struct {
		LoginType constants.LoginType `json:"login_type" binding:"required"`
		Email     string              `json:"email" binding:"required"`
		Password  string              `json:"password" binding:"required"`
	}

	ListUser struct {
		dto.PaginationRequest
	}

	RefreshToken struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
)

func (r CreateUser) ToDomain() domain.User {
	return domain.User{
		Name:   r.Name,
		Email:  r.Email,
		Roles:  r.Role,
		Status: constants.UserStatusActive,
	}
}
