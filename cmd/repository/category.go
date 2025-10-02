package repository

import (
	"context"
	"database/sql"
	default_err "errors"
	"fmt"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/errors"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, data *domain.Category) error
	CreateCategoryReturnID(ctx context.Context, data *domain.Category) (string, error)
	ListCategory(ctx context.Context, req requests.ListCategory) ([]domain.Category, int, error)
	UpdateCategory(ctx context.Context, data *domain.Category) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategory(ctx context.Context, id string) (res domain.Category, err error)
	SlugExists(ctx context.Context, slug string) (bool, error)
	GetCategoryByName(ctx context.Context, name string) (res *domain.Category, err error)
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
func (r *categoryRepository) CreateCategoryReturnID(ctx context.Context, data *domain.Category) (string, error) {
	var id string
	err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Scan(ctx, &id)
	if err != nil {
		return "", errors.CheckUniqueViolation(err)
	}
	return id, err
}

func (r *categoryRepository) ListCategory(ctx context.Context, req requests.ListCategory) ([]domain.Category, int, error) {
	var res []domain.Category
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("CreatedBy").
		Relation("EditedBy")

	if req.Search != "" {
		q.Where("category.name ILIKE ? ",
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
		ExcludeColumn("created_at", "created_by_id").
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
		Relation("CreatedBy").
		Relation("EditedBy").
		Relation("BlogArtikels").
		Where(`"category"."id" = ?`, id).Scan(ctx)
	return res, err
}
func (r *categoryRepository) GetCategoryByName(ctx context.Context, name string) (*domain.Category, error) {
	var res domain.Category
	err := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("CreatedBy").
		Relation("EditedBy").
		Relation("BlogArtikels").
		Where(`"category"."name" = ?`, name).
		Scan(ctx)
	if err != nil {
		if default_err.Is(err, sql.ErrNoRows) {
			return nil, nil // category not found
		}
		return nil, err
	}
	return &res, nil
}

func (t *categoryRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	return t.db.InitQuery(ctx).
		NewSelect().
		Model((*domain.Category)(nil)).
		Where("slug = ?", slug).
		Exists(ctx)
}
