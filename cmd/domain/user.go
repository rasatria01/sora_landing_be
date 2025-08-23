package domain

import (
	"context"
	"sora_landing_be/cmd/constants"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel
	ID             string               `bun:",pk"`
	Name           string               `bun:"name"`
	Email          string               `bun:",nullzero"`
	Roles          []constants.UserRole `bun:",array"`
	Status         constants.UserStatus `bun:","`
	Authentication *Authentication      `bun:"rel:has-one,join:id=user_id"`
	CreatedAt      time.Time            `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt      time.Time            `bun:",nullzero,notnull,default:current_timestamp"`
}

func (m *User) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.ID = ksuid.New().String()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
