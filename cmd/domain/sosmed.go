package domain

import "github.com/uptrace/bun"

type SocialMedia struct {
	bun.BaseModel
	ID       string
	Name     string
	UserName string
	Visible  bool
}
