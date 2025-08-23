package domain

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/uptrace/bun"
)

type Authentication struct {
	bun.BaseModel
	ID             string    `bun:",pk"`
	UserID         string    `bun:",nullzero"`
	User           User      `bun:"rel:belongs-to,join:user_id=id"`
	Password       string    `bun:",nullzero"`
	RefreshTokenID *string   `bun:",nullzero"`
	CreatedAt      time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt      time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (m *Authentication) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.ID = ksuid.New().String()
	}
	return nil
}
