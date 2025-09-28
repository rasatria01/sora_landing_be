package requests

import (
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/dto"
)

type TagRequest struct {
	Name string `json:"name" validate:"required"`
}

type ListTag struct {
	dto.PaginationRequest
	Search string `form:"search,omitempty"`
}

func (r *TagRequest) ToDomain(slug string) domain.Tag {
	return domain.Tag{
		Name: r.Name,
		Slug: slug,
	}
}
