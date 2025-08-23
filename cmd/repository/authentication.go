package repository

import (
	"context"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/pkg/database"

	"github.com/uptrace/bun"
)

type AuthRepository interface {
	CreateAuth(ctx context.Context, authentication *domain.Authentication) error
	GetByEmail(ctx context.Context, email string) (domain.Authentication, error)
	GetByID(ctx context.Context, id, ksuid *string) (domain.Authentication, error)
	UpdateAuth(ctx context.Context, data *domain.Authentication) error
}

type authRepository struct {
	db *database.Database
}

func NewAuthRepository(db *database.Database) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) CreateAuth(ctx context.Context, data *domain.Authentication) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Exec(ctx)
	return err
}

func (r *authRepository) UpdateAuth(ctx context.Context, data *domain.Authentication) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().Model(data).Where("id = ?", data.ID).Returning("id").Exec(ctx)
	return err
}

func (r *authRepository) GetByEmail(ctx context.Context, email string) (domain.Authentication, error) {
	var data domain.Authentication
	err := r.db.InitQuery(ctx).NewSelect().
		Model(&data).
		Where("email = ?", email).
		Relation("User", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Column("email", "status", "roles")
		}).
		Scan(ctx)
	return data, err
}

func (r *authRepository) GetByID(ctx context.Context, id, ksuid *string) (domain.Authentication, error) {
	var data domain.Authentication
	q := r.db.InitQuery(ctx).NewSelect().
		Model(&data)

	if id != nil {
		q.Where("id = ?", id)

	}

	if ksuid != nil {
		q.Where("refresh_token_id = ?", ksuid)
	}

	err := q.Scan(ctx)
	return data, err
}
