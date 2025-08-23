package requests

import (
	"sora_landing_be/cmd/constants"
	"sora_landing_be/cmd/domain"

	"github.com/segmentio/ksuid"
)

type UserAuth struct {
	AuthID         string               `json:"auth_id"`
	UserID         string               `json:"user_id"`
	Email          string               `json:"email,omitempty"`
	Role           []constants.UserRole `json:"role"`
	RefreshTokenID string               `json:"refresh_token_id,omitempty"`
}

type CreateAuth struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

func (receiver CreateAuth) ToDomain() domain.Authentication {
	return domain.Authentication{
		UserID:   receiver.UserID,
		Password: receiver.Password,
	}
}

func ToTokenPayload(record domain.Authentication) UserAuth {
	return UserAuth{
		AuthID:         record.ID,
		UserID:         record.UserID,
		Email:          record.User.Email,
		Role:           record.User.Roles,
		RefreshTokenID: ksuid.New().String(),
	}
}
