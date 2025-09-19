package response

import (
	"sora_landing_be/cmd/constants"
	"sora_landing_be/cmd/domain"
	"time"
)

type (
	// BlogArticle represents the full article response
	BlogArticle struct {
		ID          string                  `json:"id"`
		Title       string                  `json:"title"`
		Slug        string                  `json:"slug"`
		Excerpt     string                  `json:"excerpt"`
		Content     string                  `json:"content"`
		ImageURL    string                  `json:"image_url"`
		Views       int64                   `json:"views"`
		Status      constants.ArticleStatus `json:"status"`
		PublishedAt *time.Time              `json:"published_at,omitempty"`
		Category    *CategoryResponse       `json:"category,omitempty"`
		Author      *User                   `json:"author,omitempty"`
		Tags        []Tag                   `json:"tags"`
		CreatedAt   time.Time               `json:"created_at"`
		UpdatedAt   time.Time               `json:"updated_at"`
	}

	// BlogArticleList represents a summarized version for list views
	BlogArticleList struct {
		ID          string                  `json:"id"`
		Title       string                  `json:"title"`
		Slug        string                  `json:"slug"`
		Excerpt     string                  `json:"excerpt"`
		ImageURL    string                  `json:"image_url"`
		Views       int64                   `json:"views"`
		Status      constants.ArticleStatus `json:"status"`
		PublishedAt *time.Time              `json:"published_at,omitempty"`
		Category    *CategoryResponse       `json:"category,omitempty"`
		Author      *User                   `json:"author,omitempty"`
		TagCount    int                     `json:"tag_count"`
		CreatedAt   time.Time               `json:"created_at"`
	}

	// BlogArticleStats represents article statistics
	BlogArticleStats struct {
		TotalArticles     int64 `json:"total_articles"`
		PublishedArticles int64 `json:"published_articles"`
		DraftArticles     int64 `json:"draft_articles"`
		TotalViews        int64 `json:"total_views"`
	}
)

// FromDomain converts a domain BlogArtikel to a response BlogArticle
func (b *BlogArticle) FromDomain(article *domain.BlogArtikel) {
	b.ID = article.ID
	b.Title = article.Title
	b.Slug = article.Slug
	b.Excerpt = article.Excerpt
	b.Content = article.Content
	b.ImageURL = article.ImageURL
	b.Views = article.Views
	b.Status = article.Status
	b.CreatedAt = article.CreatedAt
	b.UpdatedAt = article.UpdatedAt

	if !article.PublishedAt.IsZero() {
		b.PublishedAt = &article.PublishedAt
	}

	if article.Category != nil {
		b.Category = &CategoryResponse{
			ID:   article.Category.ID,
			Name: article.Category.Name,
			Slug: article.Category.Slug,
		}
	}

	if article.Author != nil {
		b.Author = &User{
			ID:   article.Author.ID,
			Name: article.Author.Name,
		}
	}

	b.Tags = make([]Tag, len(article.Tags))
	for i, tag := range article.Tags {
		b.Tags[i] = Tag{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		}
	}
}

// FromDomain converts a domain BlogArtikel to a response BlogArticleList
func (b *BlogArticleList) FromDomain(article *domain.BlogArtikel) {
	b.ID = article.ID
	b.Title = article.Title
	b.Slug = article.Slug
	b.Excerpt = article.Excerpt
	b.ImageURL = article.ImageURL
	b.Views = article.Views
	b.Status = article.Status
	b.CreatedAt = article.CreatedAt

	if !article.PublishedAt.IsZero() {
		b.PublishedAt = &article.PublishedAt
	}

	if article.Category != nil {
		b.Category = &CategoryResponse{
			ID:   article.Category.ID,
			Name: article.Category.Name,
			Slug: article.Category.Slug,
		}
	}

	if article.Author != nil {
		b.Author = &User{
			ID:   article.Author.ID,
			Name: article.Author.Name,
		}
	}

	b.TagCount = len(article.Tags)
}
