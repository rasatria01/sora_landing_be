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
	FromURL struct {
		URL string `json:"url" validate:"required"`
	}
	UpdateArtikel struct {
		Title      string                   `json:"title" validate:"required,min=3,max=255"`
		Content    string                   `json:"content" validate:"required"`
		Excerpt    string                   `json:"excerpt" validate:"omitempty,max=500"`
		ImageURL   *string                  `json:"image_url" validate:"omitempty,url"`
		CategoryID *string                  `json:"category_id" validate:"omitempty"`
		TagIDs     []string                 `json:"tag_ids" validate:"dive,omitempty"`
		Status     *constants.ArticleStatus `json:"status" validate:"omitempty,oneof=draft published scheduled archived"`
		PublishAt  *time.Time               `json:"publish_at,omitempty" validate:"required_if=Status scheduled"`
	}

	// ListArtikel is used for querying blog articles with filters
	ListArtikel struct {
		dto.PaginationRequest
		CategoryID string                  `form:"category_id,omitempty"`
		TagID      string                  `form:"tag_id,omitempty"`
		Status     constants.ArticleStatus `form:"status,omitempty" validate:"omitempty,oneof=draft published scheduled archived"`
		Search     string                  `form:"search,omitempty"`
		StartDate  *time.Time              `form:"start_date,omitempty"`
		EndDate    *time.Time              `form:"end_date,omitempty"`
		SortBy     string                  `form:"sort_by,omitempty" validate:"omitempty,oneof=created_at published_at views title"`
		SortOrder  string                  `form:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
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
func (r *UpdateArtikel) ToDomain(userID string, slug string) *domain.BlogArtikel {
	article := &domain.BlogArtikel{
		Slug:     slug,
		AuthorID: userID,
		Tags:     make([]*domain.Tag, 0), // tags handled in service
	}

	if r.Title != "" {
		article.Title = r.Title
	}

	if r.Content != "" {
		article.Content = r.Content
	}

	if r.Excerpt != "" {
		article.Excerpt = r.Excerpt
	}

	if r.ImageURL != nil {
		article.ImageURL = *r.ImageURL
	}

	if r.Status != nil {
		article.Status = *r.Status

		// handle published/scheduled
		if *r.Status == constants.StatusPublished {
			article.PublishedAt = time.Now()
		} else if *r.Status == constants.StatusScheduled && r.PublishAt != nil {
			article.PublishedAt = *r.PublishAt
		}
	}

	if r.CategoryID != nil {
		article.CategoryID = *r.CategoryID
	}

	return article
}
