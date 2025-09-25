package response

import "sora_landing_be/cmd/domain"

type (
	Tag struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Slug      string `json:"slug"`
		CreatedBy string `json:"created_by"`
		UpdatedBy string `json:"updated_by"`
	}
)

func NewListTag(tags []domain.Tag) []Tag {
	var res []Tag
	for _, tag := range tags {
		res = append(res, Tag{
			ID:        tag.ID,
			Name:      tag.Name,
			Slug:      tag.Slug,
			CreatedBy: tag.CreatedBy.Name,
			UpdatedBy: tag.EditedBy.Name,
		})
	}

	return res
}
func NewListTags(tags []*domain.Tag) []*Tag {
	var res []*Tag
	for _, tag := range tags {
		res = append(res, &Tag{
			ID:        tag.ID,
			Name:      tag.Name,
			Slug:      tag.Slug,
			CreatedBy: tag.CreatedBy.Name,
			UpdatedBy: tag.EditedBy.Name,
		})
	}

	return res
}

func NewTag(tag domain.Tag) Tag {
	return Tag{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedBy: tag.CreatedBy.Name,
		UpdatedBy: tag.EditedBy.Name,
	}

}
