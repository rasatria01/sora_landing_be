package requests

import "sora_landing_be/cmd/dto"

// Public blog requests can have more limited fields compared to admin
type PublicListArticle struct {
	dto.PaginationRequest
	CategoryID string `json:"category_id" form:"category_id"`
	TagID      string `json:"tag_id" form:"tag_id"`
	Search     string `json:"search" form:"search"`
	SortBy     string `json:"sort_by" form:"sort_by" binding:"omitempty,oneof=published_at views title"`
	SortOrder  string `json:"sort_order" form:"sort_order" binding:"omitempty,oneof=asc desc"`
}
