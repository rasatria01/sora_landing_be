package domain

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/uptrace/bun"
)

type BaseEntity struct {
	ID        string    `bun:",pk"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

func (m *BaseEntity) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.ID = ksuid.New().String()
		if m.CreatedAt.IsZero() {
			m.CreatedAt = time.Now()
		}
		m.UpdatedAt = m.CreatedAt
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
