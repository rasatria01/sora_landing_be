package domain

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/uptrace/bun"
)

type FileUpload struct {
	bun.BaseModel `bun:"table:file_uploads,alias:fu"`
	ID            string    `bun:"id,pk,type:VARCHAR(27)"`                 // KSUID primary key
	FileName      string    `bun:"file_name,notnull"`                      // Original file name
	FilePath      string    `bun:"file_path,notnull"`                      // Storage path
	FileSize      int64     `bun:"file_size,notnull"`                      // File size in bytes
	ContentType   string    `bun:"content_type,notnull"`                   // MIME type
	Module        string    `bun:"module,nullzero"`                        // Optional: module/feature using this file
	ReferenceID   string    `bun:"reference_id,nullzero,type:VARCHAR(27)"` // Optional: ID of related entity
	AuthorID      string    `bun:"uploaded_by,notnull,type:VARCHAR(27)"`
	Author        *User     `bun:"rel:belongs-to,join:uploaded_by=id"` // User ID who uploaded
	IsPublic      bool      `bun:"is_public,notnull"`                  // Whether file is publicly accessible
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt     time.Time `bun:"deleted_at,soft_delete,nullzero"`
}

func (m *FileUpload) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.ID = ksuid.New().String()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
