package requests

import (
	"sora_landing_be/cmd/constants"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/dto"
	"time"
)

type (
	// BlogArtikel is used for creating and updating blog articles
	BlogArtikel struct {
		Title      string                  `json:"title" validate:"required,min=3,max=255"`
		Content    string                  `json:"content" validate:"required"`
		Excerpt    string                  `json:"excerpt" validate:"omitempty,max=500"`
		ImageURL   string                  `json:"image_url" validate:"omitempty,url"`
		CategoryID string                  `json:"category_id" validate:"required"`
		TagIDs     []string                `json:"tag_ids" validate:"dive,required"`
		Status     constants.ArticleStatus `json:"status" validate:"required,oneof=draft published scheduled archived"`
		PublishAt  *time.Time              `json:"publish_at,omitempty" validate:"required_if=Status scheduled"`
	}

	// ListArtikel is used for querying blog articles with filters
	ListArtikel struct {
		dto.PaginationRequest
		CategoryID string                  `json:"category_id,omitempty"`
		TagID      string                  `json:"tag_id,omitempty"`
		Status     constants.ArticleStatus `json:"status,omitempty" validate:"omitempty,oneof=draft published scheduled archived"`
		Search     string                  `json:"search,omitempty"`
		StartDate  *time.Time              `json:"start_date,omitempty"`
		EndDate    *time.Time              `json:"end_date,omitempty"`
		SortBy     string                  `json:"sort_by,omitempty" validate:"omitempty,oneof=created_at published_at views title"`
		SortOrder  string                  `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
	}

	// UpdateArticleStatus is used for changing article status
	UpdateArticleStatus struct {
		Status    constants.ArticleStatus `json:"status" validate:"required,oneof=draft published scheduled archived"`
		PublishAt *time.Time              `json:"publish_at,omitempty" validate:"required_if=Status scheduled"`
	}
)

func (r *BlogArtikel) ToDomain(userID string, slug string) *domain.BlogArtikel {
	article := &domain.BlogArtikel{
		Title:      r.Title,
		Slug:       slug,
		Content:    r.Content,
		Excerpt:    r.Excerpt,
		ImageURL:   r.ImageURL,
		Status:     r.Status,
		CategoryID: r.CategoryID,
		AuthorID:   userID,
		Tags:       make([]*domain.Tag, 0), // will be filled later by service
	}

	if r.Status == constants.StatusPublished {
		now := time.Now()
		article.PublishedAt = now
	} else if r.Status == constants.StatusScheduled && r.PublishAt != nil {
		article.PublishedAt = *r.PublishAt
	}

	return article
}
