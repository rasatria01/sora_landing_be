package dto

// BlogStats represents statistics about blog articles
type BlogStats struct {
	TotalArticles     int64 `json:"total_articles"`
	PublishedArticles int64 `json:"published_articles"`
	DraftArticles     int64 `json:"draft_articles"`
	TotalViews        int64 `json:"total_views"`
}
