package response

import (
	"sora_landing_be/cmd/domain"
	"time"
)

type CategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy *string   `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToCategoryResponse(category domain.Category) CategoryResponse {
	var updatedBy *string
	if category.EditedBy != nil {
		updatedBy = &category.EditedBy.Name
	}
	return CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedBy: category.CreatedBy.Name,
		UpdatedBy: updatedBy,
		UpdatedAt: category.UpdatedAt,
	}
}

func NewListCat(tags []domain.Category) []CategoryResponse {
	var res []CategoryResponse
	for _, tag := range tags {
		var updatedBy *string
		if tag.EditedBy != nil {
			updatedBy = &tag.EditedBy.Name
		}
		res = append(res, CategoryResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			Slug:      tag.Slug,
			CreatedBy: tag.CreatedBy.Name,
			UpdatedBy: updatedBy,
			UpdatedAt: tag.UpdatedAt,
		})
	}

	return res
}
