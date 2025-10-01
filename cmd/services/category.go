package services

import (
	"context"
	"database/sql"
	"sora_landing_be/cmd/dto"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/cmd/dto/response"
	"sora_landing_be/cmd/repository"
	"sora_landing_be/pkg/authentication"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/utils"

	"github.com/uptrace/bun"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, payload requests.Category) error
	ListCategory(ctx context.Context, payload requests.ListCategory) (dto.PaginationResponse[response.CategoryResponse], error)
	UpdateCategory(ctx context.Context, id string, payload requests.Category) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategory(ctx context.Context, id string) (response.CategoryResponse, error)
}

type catService struct {
	catRepo repository.CategoryRepository
}

func NewCatService(catRepo repository.CategoryRepository) CategoryService {
	return &catService{
		catRepo: catRepo,
	}
}

func (t *catService) CreateCategory(ctx context.Context, payload requests.Category) error {
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		uniqueSlug, err := utils.GenerateUniqueSlug(ctx, t.catRepo, payload.Name)
		if err != nil {
			return err
		}
		data := payload.ToDomain(uniqueSlug)
		data.CreatedByID = authentication.GetUserDataFromToken(ctx).UserID
		err = t.catRepo.CreateCategory(ctx, &data)
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

func (a *catService) ListCategory(ctx context.Context, payload requests.ListCategory) (dto.PaginationResponse[response.CategoryResponse], error) {
	var paginateRes dto.PaginationResponse[response.CategoryResponse]
	res, count, err := a.catRepo.ListCategory(ctx, payload)
	if err != nil {
		return paginateRes, err
	}

	paginateRes = dto.NewPaginationResponse(payload.PaginationRequest, count, response.NewListCat(res))
	return paginateRes, nil
}

func (a *catService) UpdateCategory(ctx context.Context, id string, payload requests.Category) error {
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		uniqueSlug, err := utils.GenerateUniqueSlug(ctx, a.catRepo, payload.Name)
		if err != nil {
			return err
		}
		data := payload.ToDomain(uniqueSlug)
		data.ID = id
		edited := authentication.GetUserDataFromToken(ctx).UserID
		data.EditedByID = &edited

		err = a.catRepo.UpdateCategory(ctx, &data)
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

func (a *catService) DeleteCategory(ctx context.Context, id string) error {

	err := a.catRepo.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *catService) GetCategory(ctx context.Context, id string) (response.CategoryResponse, error) {
	var res response.CategoryResponse
	data, err := a.catRepo.GetCategory(ctx, id)
	if err != nil {
		return res, err
	}

	res = response.ToCategoryResponse(data)
	return res, nil
}
