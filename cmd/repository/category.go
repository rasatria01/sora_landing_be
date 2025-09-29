package repository

import (
	"context"
	"fmt"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/errors"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, data *domain.Category) error
	ListCategory(ctx context.Context, req requests.ListCategory) ([]domain.Category, int, error)
	UpdateCategory(ctx context.Context, data *domain.Category) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategory(ctx context.Context, id string) (res domain.Category, err error)
	SlugExists(ctx context.Context, slug string) (bool, error)
}

type categoryRepository struct {
	db *database.Database
}

func NewCatRepository(db *database.Database) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) CreateCategory(ctx context.Context, data *domain.Category) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Exec(ctx)
	if err != nil {
		return errors.CheckUniqueViolation(err)
	}
	return err
}

func (r *categoryRepository) ListCategory(ctx context.Context, req requests.ListCategory) ([]domain.Category, int, error) {
	var res []domain.Category
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res)

	if req.Search != "" {
		q.Where("name ILIKE ? ",
			fmt.Sprintf("%%%s%%", req.Search))
	}
	q.Limit(req.PageSize).
		Offset(req.CalculateOffset()).
		Order(fmt.Sprintf("%s %s", req.OrderBy, req.OrderDir))
	total, err := q.ScanAndCount(ctx)
	return res, total, err
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, data *domain.Category) error {
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model(data).
		Where("id = ?", data.ID).
		ExcludeColumn("created_at").
		Returning("id").
		Exec(ctx)
	return err
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).
		NewDelete().
		Model((*domain.Category)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *categoryRepository) GetCategory(ctx context.Context, id string) (res domain.Category, err error) {
	err = r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("BlogArtikels").
		Where(`"category"."id" = ?`, id).Scan(ctx)
	return res, err
}

func (t *categoryRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	return t.db.InitQuery(ctx).
		NewSelect().
		Model((*domain.Category)(nil)).
		Where("slug = ?", slug).
		Exists(ctx)
}
