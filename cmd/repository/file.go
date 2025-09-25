package repository

import (
	"context"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/pkg/database"
)

type FileRepository interface {
	CreateFile(ctx context.Context, file *domain.FileUpload) error
	GetFileByID(ctx context.Context, id string) (domain.FileUpload, error)
	DeleteFile(ctx context.Context, id string) error
	GetPublicFileByID(ctx context.Context, id string) (domain.FileUpload, error)
}

type fileRepository struct {
	db *database.Database
}

func NewFileRepository(db *database.Database) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) CreateFile(ctx context.Context, file *domain.FileUpload) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(file).Exec(ctx)
	return err
}

func (r *fileRepository) GetFileByID(ctx context.Context, id string) (domain.FileUpload, error) {
	var file domain.FileUpload
	err := r.db.InitQuery(ctx).NewSelect().Model(&file).Where("id = ?", id).Scan(ctx)
	return file, err
}

func (r *fileRepository) GetPublicFileByID(ctx context.Context, id string) (domain.FileUpload, error) {
	var file domain.FileUpload
	err := r.db.InitQuery(ctx).NewSelect().Model(&file).
		Where("id = ? AND is_public = true", id).
		Scan(ctx)
	return file, err
}

func (r *fileRepository) DeleteFile(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Table("file_uploads").
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Exec(ctx)
	return err
}
