package repository

import (
	"context"
	"fmt"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/errors"
)

type UserRepository interface {
	CreateUser(ctx context.Context, data *domain.User) error
	ListUser(ctx context.Context, req requests.ListUser) ([]domain.User, int, error)
	UpdateUser(ctx context.Context, data *domain.User) error
	DeleteUser(ctx context.Context, id string) error
	GetUser(ctx context.Context, id string) (res domain.User, err error)
}

type userRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, data *domain.User) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Exec(ctx)
	if err != nil {
		return errors.CheckUniqueViolation(err)
	}
	return err
}

func (r *userRepository) ListUser(ctx context.Context, req requests.ListUser) ([]domain.User, int, error) {
	var res []domain.User
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

func (r *userRepository) UpdateUser(ctx context.Context, data *domain.User) error {
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model(data).
		Where("id = ?", data.ID).
		ExcludeColumn("created_at").
		Returning("id").
		Exec(ctx)
	return err
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).
		NewDelete().
		Model((*domain.User)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *userRepository) GetUser(ctx context.Context, id string) (res domain.User, err error) {
	err = r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Where(`"user"."id" = ?`, id).Scan(ctx)
	return res, err
}
