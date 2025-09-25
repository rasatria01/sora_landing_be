package response

import (
	"sora_landing_be/cmd/domain"
	"time"
)

type FileUploadResponse struct {
	ID          string    `json:"id"`
	FileName    string    `json:"file_name"`
	FileURL     string    `json:"file_url"`
	FileSize    int64     `json:"file_size"`
	ContentType string    `json:"content_type"`
	Module      string    `json:"module,omitempty"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

func NewFileUploadResponse(file *domain.FileUpload, baseURL string) FileUploadResponse {
	return FileUploadResponse{
		ID:          file.ID,
		FileName:    file.FileName,
		FileURL:     baseURL + "/v1/public/files/" + file.ID,
		FileSize:    file.FileSize,
		ContentType: file.ContentType,
		Module:      file.Module,
		UploadedAt:  file.CreatedAt,
	}
}
