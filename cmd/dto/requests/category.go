package requests

import (
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/dto"
)

type (
	Category struct {
		Name string `json:"name" validate:"required"`
	}
	ListCategory struct {
		dto.PaginationRequest
	}
)

func (c *Category) ToDomain(slug string) domain.Category {
	return domain.Category{
		Name: c.Name,
		Slug: slug,
	}
}
