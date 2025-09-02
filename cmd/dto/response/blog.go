package response

import "time"

type (
	BlogArticle struct {
		ID          string            `json:"id"`
		Title       string            `json:"title"`
		Slug        string            `json:"slug"`
		Excerpt     string            `json:"excerpt"`
		Content     string            `json:"content"`
		ImageURL    string            `json:"image_url"`
		Views       int               `json:"views"`
		Status      string            `json:"status"`
		PublishedAt *time.Time        `json:"published_at"`
		Category    *CategoryResponse `json:"category,omitempty"`
		Tags        []Tag             `json:"tags"`
		CreatedAt   time.Time         `json:"created_at"`
		UpdatedAt   time.Time         `json:"updated_at"`
	}
)
