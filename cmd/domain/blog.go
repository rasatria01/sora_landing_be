package domain

import (
	"sora_landing_be/cmd/constants"
	"time"

	"github.com/uptrace/bun"
)

type BlogArtikel struct {
	bun.BaseModel `bun:"table:blog_artikels,alias:ba"`
	BaseEntity

	Title       string                  `bun:",notnull"`
	Slug        string                  `bun:",unique,notnull"`
	Content     string                  `bun:",type:text,notnull"`
	Excerpt     string                  `bun:",type:text"`
	ImageURL    string                  `bun:",nullzero"` // optional feature imag
	CategoryID  string                  `bun:",notnull"`
	Category    *Category               `bun:"rel:belongs-to,join:category_id=id"`
	AuthorID    string                  `bun:",notnull"`
	Author      *User                   `bun:"rel:belongs-to,join:author_id=id"`
	Status      constants.ArticleStatus `bun:",notnull,default:'draft'"` // draft, published, archived
	Views       int64                   `bun:",default:0"`
	PublishedAt time.Time               `bun:",nullzero"`
	Tags        []*Tag                  `bun:"m2m:article_tags,join:Article=Tag"`
}
