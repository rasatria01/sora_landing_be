package domain

import (
	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:categories,alias:c"`
	BaseEntity

	Name string `bun:",unique,notnull"`
	Slug string `bun:",unique,notnull"`

	// Reverse relation
	BlogArtikels []*BlogArtikel `bun:"rel:has-many,join:id=category_id"`
}
