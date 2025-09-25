package response

import (
	"sora_landing_be/cmd/domain"
	"time"
)

type CategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

func ToCategoryResponse(category domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
	}
}

func NewListCat(tags []domain.Category) []CategoryResponse {
	var res []CategoryResponse
	for _, tag := range tags {
		res = append(res, CategoryResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			Slug:      tag.Slug,
			CreatedAt: tag.CreatedAt,
		})
	}

	return res
}
