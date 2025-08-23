package dto

import "math"

type PaginationRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size,default=10"`
	OrderBy  string `form:"order_by,default=updated_at"`
	OrderDir string `form:"order_dir,default=desc" binding:"oneof=desc asc"`
}

type PaginationResponse[T any] struct {
	CurrentPage int  `json:"current_page"`
	PageSize    int  `json:"page_size"`
	TotalItems  int  `json:"total_items"`
	TotalPages  int  `json:"total_pages"`
	HasPrevious bool `json:"has_previous"`
	HasNext     bool `json:"has_next"`
	Data        []T  `json:"data"`
}

func NewPaginationResponse[T any](req PaginationRequest, totalItems int, data []T) PaginationResponse[T] {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(req.PageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationResponse[T]{
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasPrevious: req.Page > 1,
		HasNext:     req.Page < totalPages,
		Data:        data,
	}
}

func (p PaginationRequest) CalculateOffset() int {
	return (p.Page - 1) * p.PageSize
}
