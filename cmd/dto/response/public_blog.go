package response

import (
	"sora_landing_be/cmd/domain"
	"time"
)

// PublicArticleList is a simplified version of article for list views
type PublicArticleList struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Excerpt     string     `json:"excerpt"`
	ImageURL    string     `json:"image_url"`
	Views       int64      `json:"views"`
	PublishedAt *time.Time `json:"published_at"`

	// Simplified related data
	Category *CategoryResponse   `json:"category"`
	Author   *PublicAuthorDetail `json:"author"`
	Tags     []Tag               `json:"tags"`
}

// PublicArticleDetail is the full article view for public
type PublicArticleDetail struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	Excerpt     string     `json:"excerpt"`
	ImageURL    string     `json:"image_url"`
	Views       int64      `json:"views"`
	PublishedAt *time.Time `json:"published_at"`

	// Related data
	Category *CategoryResponse   `json:"category"`
	Author   *PublicAuthorDetail `json:"author"`
	Tags     []Tag               `json:"tags"`

	// Related articles
	RelatedArticles []PublicArticleList `json:"related_articles"`
}

// PublicAuthorDetail contains non-sensitive author information
type PublicAuthorDetail struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (p *PublicArticleList) FromDomain(article *domain.BlogArtikel) {
	p.ID = article.ID
	p.Title = article.Title
	p.Slug = article.Slug
	p.Excerpt = article.Excerpt
	p.ImageURL = article.ImageURL
	p.Views = article.Views
	if !article.PublishedAt.IsZero() {
		p.PublishedAt = &article.PublishedAt
	}

	if article.Category != nil {
		p.Category = &CategoryResponse{
			ID:   article.Category.ID,
			Name: article.Category.Name,
			Slug: article.Category.Slug,
		}
	}
	if article.Author != nil {
		p.Author = &PublicAuthorDetail{
			ID:   article.Author.ID,
			Name: article.Author.Name,
		}
	}
	p.Tags = make([]Tag, len(article.Tags))
	for i, tag := range article.Tags {
		p.Tags[i] = Tag{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		}
	}
}

func (p *PublicArticleDetail) FromDomain(article *domain.BlogArtikel, relatedArticles []domain.BlogArtikel) {
	p.ID = article.ID
	p.Title = article.Title
	p.Slug = article.Slug
	p.Content = article.Content
	p.Excerpt = article.Excerpt
	p.ImageURL = article.ImageURL
	p.Views = article.Views
	if !article.PublishedAt.IsZero() {
		p.PublishedAt = &article.PublishedAt
	}

	if article.Category != nil {
		p.Category = &CategoryResponse{
			ID:   article.Category.ID,
			Name: article.Category.Name,
			Slug: article.Category.Slug,
		}
	}
	if article.Author != nil {
		p.Author = &PublicAuthorDetail{
			ID:   article.Author.ID,
			Name: article.Author.Name,
		}
	}
	p.Tags = make([]Tag, len(article.Tags))
	for i, tag := range article.Tags {
		p.Tags[i] = Tag{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		}
	}

	// Convert related articles
	if len(relatedArticles) > 0 {
		p.RelatedArticles = make([]PublicArticleList, len(relatedArticles))
		for i, rel := range relatedArticles {
			p.RelatedArticles[i].FromDomain(&rel)
		}
	}
}
