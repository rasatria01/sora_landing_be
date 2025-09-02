package requests

import (
	"sora_landing_be/cmd/constants"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/dto"
)

type (
	BlogArtikel struct {
		Title      string                  `json:"title" validate:"required"`
		Content    string                  `json:"content" validate:"required"`
		ImageURL   string                  `json:"image_url"`
		CategoryID string                  `json:"category_id" validate:"required"`
		TagIDs     []string                `json:"tag_ids"`
		Status     constants.ArticleStatus `json:"status" validate:"required,oneof=draft published scheduled archived"`
	}

	ListArtikel struct {
		dto.PaginationRequest
	}
)

func (r *BlogArtikel) ToDomain(userID string, slug string, excerpt string) *domain.BlogArtikel {
	return &domain.BlogArtikel{
		Title:      r.Title,
		Slug:       slug,
		Excerpt:    excerpt,
		Content:    r.Content,
		ImageURL:   r.ImageURL,
		Status:     r.Status,
		CategoryID: r.CategoryID,
		AuthorID:   userID,
		Tags:       make([]domain.Tag, 0), // will be filled later by service
	}
}
