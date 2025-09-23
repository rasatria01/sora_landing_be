package domain

import "github.com/uptrace/bun"

type ArticleTag struct {
	bun.BaseModel `bun:"table:article_tags"`

	ArticleID string       `bun:"blog_article_id,pk"`
	Article   *BlogArtikel `bun:"rel:belongs-to,join:blog_article_id=id"`

	TagID string `bun:",pk"`
	Tag   *Tag   `bun:"rel:belongs-to,join:tag_id=id"`
}
