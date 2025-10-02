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

type TagRepository interface {
	CreateTag(ctx context.Context, data *domain.Tag) error
	CreateTagReturnID(ctx context.Context, data *domain.Tag) (string, error)
	ListTag(ctx context.Context, req requests.ListTag) ([]domain.Tag, int, error)
	UpdateTag(ctx context.Context, data *domain.Tag) error
	DeleteTag(ctx context.Context, id string) error
	GetTag(ctx context.Context, id string) (res domain.Tag, err error)
	GetTagByName(ctx context.Context, name string) (*domain.Tag, error)
	SlugExists(ctx context.Context, slug string) (bool, error)
}

type tagRepository struct {
	db *database.Database
}

func NewTagRepository(db *database.Database) TagRepository {
	return &tagRepository{
		db: db,
	}
}

func (r *tagRepository) CreateTag(ctx context.Context, data *domain.Tag) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Exec(ctx)
	if err != nil {
		return errors.CheckUniqueViolation(err)
	}
	return err
}
func (r *tagRepository) CreateTagReturnID(ctx context.Context, data *domain.Tag) (string, error) {
	var id string
	err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Scan(ctx, &id)
	if err != nil {
		return "", errors.CheckUniqueViolation(err)
	}
	return id, err
}

func (r *tagRepository) ListTag(ctx context.Context, req requests.ListTag) ([]domain.Tag, int, error) {
	var res []domain.Tag
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("CreatedBy").
		Relation("EditedBy")
	if req.Search != "" {
		q.Where("tag.name ILIKE ? ",
			fmt.Sprintf("%%%s%%", req.Search))
	}

	q.Limit(req.PageSize).
		Offset(req.CalculateOffset()).
		Order(fmt.Sprintf("%s %s", req.OrderBy, req.OrderDir))

	total, err := q.ScanAndCount(ctx)
	return res, total, err
}

func (r *tagRepository) UpdateTag(ctx context.Context, data *domain.Tag) error {
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model(data).
		Where("id = ?", data.ID).
		ExcludeColumn("created_at", "created_by_id").
		Returning("id").
		Exec(ctx)
	return err
}

func (r *tagRepository) DeleteTag(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).
		NewDelete().
		Model((*domain.Tag)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *tagRepository) GetTag(ctx context.Context, id string) (res domain.Tag, err error) {
	err = r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("CreatedBy").
		Relation("EditedBy").
		Where(`"tag"."id" = ?`, id).Scan(ctx)
	return res, err
}
func (r *tagRepository) GetTagByName(ctx context.Context, name string) (*domain.Tag, error) {
	var res domain.Tag
	err := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("CreatedBy").
		Relation("EditedBy").
		Where(`"tag"."name" = ?`, name).Scan(ctx)
	if err != nil {
		if default_err.Is(err, sql.ErrNoRows) {
			return nil, nil // category not found
		}
		return nil, err
	}
	return &res, err
}

func (t *tagRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	return t.db.InitQuery(ctx).
		NewSelect().
		Model((*domain.Tag)(nil)).
		Where("slug = ?", slug).
		Exists(ctx)
}
