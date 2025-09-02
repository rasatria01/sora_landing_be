package constants

type ArticleStatus string

const (
	StatusDraft     ArticleStatus = "draft"
	StatusPublished ArticleStatus = "published"
	StatusScheduled ArticleStatus = "scheduled"
	StatusArchived  ArticleStatus = "archived"
)
