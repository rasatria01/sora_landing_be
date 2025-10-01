package response

import (
	"sora_landing_be/cmd/domain"
	"time"
)

type (
	User struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	Profile struct {
		ID     string   `json:"id"`
		Name   string   `json:"name"`
		Email  string   `json:"email"`
		Permit []string `json:"permit"`
	}
)

func NewProfile(user domain.User, permit []string) Profile {
	return Profile{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Permit: permit,
	}
}

func NewListUser(users []domain.User) []User {
	var res []User
	for _, user := range users {
		res = append(res, User{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return res
}

func NewUser(user domain.User) User {
	return User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

}
