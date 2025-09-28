package response

import (
	"sora_landing_be/cmd/domain"
	"time"
)

type (
	Tag struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Slug      string    `json:"slug"`
		CreatedBy string    `json:"created_by"`
		UpdatedBy *string   `json:"updated_by"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func NewListTag(tags []domain.Tag) []Tag {
	var res []Tag
	for _, tag := range tags {
		var updatedBy *string
		if tag.EditedBy != nil {
			updatedBy = &tag.EditedBy.Name
		}
		res = append(res, Tag{
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
func NewListTags(tags []*domain.Tag) []*Tag {
	var res []*Tag
	for _, tag := range tags {
		var updatedBy *string
		if tag.EditedBy != nil {
			updatedBy = &tag.EditedBy.Name
		}
		res = append(res, &Tag{
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

func NewTag(tag domain.Tag) Tag {
	var updatedBy *string
	if tag.EditedBy != nil {
		updatedBy = &tag.EditedBy.Name
	}
	return Tag{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedBy: tag.CreatedBy.Name,
		UpdatedBy: updatedBy,
		UpdatedAt: tag.UpdatedAt,
	}

}
