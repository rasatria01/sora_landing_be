package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type DemoEntry struct {
	bun.BaseModel `bun:"table:demo"`
	BaseEntity
	Nama    string    `bun:",notnull"`
	Brand   string    `bun:",notnull"`
	NoHP    string    `bun:",notnull"`
	Email   string    `bun:",nullzero"`
	Waktu   string    `bun:",notnull"`
	Tanggal time.Time `bun:",notnull"`
}
