package repository

import (
	"context"
	"fmt"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/errors"
)

type DemoRepository interface {
	CreateDemo(ctx context.Context, entry *domain.DemoEntry) error
	GetDemoByID(ctx context.Context, id string) (domain.DemoEntry, error)
	ListDemos(ctx context.Context, req requests.ListDemo) ([]domain.DemoEntry, int, error)
	UpdateDemo(ctx context.Context, entry *domain.DemoEntry) error
	DeleteDemo(ctx context.Context, id string) error
	ExportDemo(ctx context.Context, req requests.ExportDemo) ([]domain.DemoEntry, error)
}

type demoRepository struct {
	db *database.Database
}

func NewDemoRepository(db *database.Database) DemoRepository {
	return &demoRepository{db: db}
}

func (r *demoRepository) CreateDemo(ctx context.Context, entry *domain.DemoEntry) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(entry).Returning("id").Exec(ctx)
	if err != nil {
		return errors.CheckUniqueViolation(err)
	}

	return err
}

func (r *demoRepository) GetDemoByID(ctx context.Context, id string) (domain.DemoEntry, error) {
	var entry domain.DemoEntry
	err := r.db.
		InitQuery(ctx).
		NewSelect().
		Model(&entry).
		Where("id = ?", id).
		Scan(ctx)
	return entry, err
}

func (r *demoRepository) ListDemos(ctx context.Context, req requests.ListDemo) ([]domain.DemoEntry, int, error) {
	var entries []domain.DemoEntry
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&entries)
	if req.Search != "" {
		q.Where("nama ILIKE ? ",
			fmt.Sprintf("%%%s%%", req.Search))
	}
	q.Limit(req.PageSize).
		Offset(req.CalculateOffset()).
		Order(fmt.Sprintf("%s %s", req.OrderBy, req.OrderDir))

	total, err := q.ScanAndCount(ctx)
	return entries, total, err
}

func (r *demoRepository) UpdateDemo(ctx context.Context, entry *domain.DemoEntry) error {
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model(entry).
		Where("id = ?", entry.ID).
		ExcludeColumn("created_at").
		Returning("id").
		Exec(ctx)
	return err
}

func (r *demoRepository) DeleteDemo(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).
		NewDelete().
		Model((*domain.DemoEntry)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *demoRepository) ExportDemo(ctx context.Context, req requests.ExportDemo) ([]domain.DemoEntry, error) {
	var entries []domain.DemoEntry
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&entries)
	if req.StartDate != nil {
		q.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != nil {
		q.Where("created_at <= ?", req.EndDate)
	}
	err := q.Scan(ctx)
	return entries, err
}
