package services

import (
	"context"
	"database/sql"
	"sora_landing_be/cmd/dto"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/cmd/dto/response"
	"sora_landing_be/cmd/repository"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/utils"

	"github.com/uptrace/bun"
)

type TagService interface {
	CreateTag(ctx context.Context, payload requests.TagRequest) error
	ListTag(ctx context.Context, payload requests.ListTag) (dto.PaginationResponse[response.Tag], error)
	UpdateTag(ctx context.Context, id string, payload requests.TagRequest) error
	DeleteTag(ctx context.Context, id string) error
	GetTag(ctx context.Context, id string) (response.Tag, error)
}

type tagService struct {
	tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{
		tagRepo: tagRepo,
	}
}

func (t *tagService) CreateTag(ctx context.Context, payload requests.TagRequest) error {
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		uniqueSlug, err := utils.GenerateUniqueSlug(ctx, t.tagRepo, payload.Name)
		if err != nil {
			return err
		}
		data := payload.ToDomain(uniqueSlug)

		err = t.tagRepo.CreateTag(ctx, &data)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *tagService) ListTag(ctx context.Context, payload requests.ListTag) (dto.PaginationResponse[response.Tag], error) {
	var paginateRes dto.PaginationResponse[response.Tag]
	res, count, err := a.tagRepo.ListTag(ctx, payload)
	if err != nil {
		return paginateRes, err
	}

	paginateRes = dto.NewPaginationResponse(payload.PaginationRequest, count, response.NewListTag(res))
	return paginateRes, nil
}

func (a *tagService) UpdateTag(ctx context.Context, id string, payload requests.TagRequest) error {
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		uniqueSlug, err := utils.GenerateUniqueSlug(ctx, a.tagRepo, payload.Name)
		if err != nil {
			return err
		}
		data := payload.ToDomain(uniqueSlug)
		data.ID = id

		err = a.tagRepo.UpdateTag(ctx, &data)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *tagService) DeleteTag(ctx context.Context, id string) error {

	err := a.tagRepo.DeleteTag(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *tagService) GetTag(ctx context.Context, id string) (response.Tag, error) {
	var res response.Tag
	data, err := a.tagRepo.GetTag(ctx, id)
	if err != nil {
		return res, err
	}

	res = response.NewTag(data)
	return res, nil
}
