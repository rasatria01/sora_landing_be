package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/dto/response"
	"sora_landing_be/cmd/repository"
	"sora_landing_be/pkg/authentication"
	"sora_landing_be/pkg/errors"
	"sora_landing_be/pkg/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type FileService interface {
	UploadFiles(ctx context.Context, files []*multipart.FileHeader, module string) ([]response.FileUploadResponse, error)
	GetPublicFile(ctx context.Context, id string) (domain.FileUpload, error)
	DeleteFile(ctx context.Context, id string) error
}

type fileService struct {
	fileRepo repository.FileRepository
	baseURL  string
}

func NewFileService(fileRepo repository.FileRepository, baseURL string) FileService {
	return &fileService{
		fileRepo: fileRepo,
		baseURL:  baseURL,
	}
}

func (s *fileService) UploadFiles(ctx context.Context, files []*multipart.FileHeader, module string) ([]response.FileUploadResponse, error) {
	responses := make([]response.FileUploadResponse, 0, len(files))

	for _, file := range files {
		// Check file size (e.g., 10MB limit)
		if file.Size > 10*1024*1024 {
			return nil, errors.NewValidationError("File size exceeds limit",
				map[string]string{"file": fmt.Sprintf("File %s exceeds 10MB limit", file.Filename)})
		}

		// Check file type
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !isAllowedFileType(ext) {
			return nil, errors.NewValidationError("Invalid file type",
				map[string]string{"file": fmt.Sprintf("File type %s not allowed", ext)})
		}

		// Process and save file
		var filePath string
		var fileSize int64
		var contentType string

		if utils.IsImage(file.Filename) {
			// Process image file
			processed, err := utils.ProcessImage(file)
			if err != nil {
				return nil, fmt.Errorf("failed to process image: %w", err)
			}

			// Generate unique filename
			ext := filepath.Ext(file.Filename)
			uniqueName := fmt.Sprintf("%s_%d%s", strings.TrimSuffix(file.Filename, ext), time.Now().UnixNano(), ext)
			filePath = fmt.Sprintf("uploads/%s/%s", module, uniqueName)

			// Create directory if not exists
			err = os.MkdirAll(filepath.Dir(filePath), 0755)
			if err != nil {
				return nil, fmt.Errorf("failed to create directory: %w", err)
			}

			// Save processed image
			err = os.WriteFile(filePath, processed, 0644)
			if err != nil {
				return nil, fmt.Errorf("failed to save processed file: %w", err)
			}

			fileSize = int64(len(processed))
			contentType = "image/" + utils.GetImageFormat(file.Filename)
		} else {
			// Save non-image file as is
			filePath = fmt.Sprintf("uploads/%s/%s", module, file.Filename)
			err := ctx.(*gin.Context).SaveUploadedFile(file, filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to save file: %w", err)
			}
			fileSize = file.Size
			contentType = file.Header.Get("Content-Type")
		}
		userID := authentication.GetUserDataFromToken(ctx).UserID

		// Create file record
		fileUpload := &domain.FileUpload{
			FileName:    file.Filename,
			FilePath:    filePath,
			FileSize:    fileSize,
			ContentType: contentType,
			Module:      module,
			AuthorID:    userID,
			IsPublic:    true, // Since this is for public endpoint
		}

		if err := s.fileRepo.CreateFile(ctx, fileUpload); err != nil {
			return nil, fmt.Errorf("failed to create file record: %w", err)
		}

		responses = append(responses, response.NewFileUploadResponse(fileUpload, s.baseURL))
	}

	return responses, nil
}

func (s *fileService) GetPublicFile(ctx context.Context, id string) (domain.FileUpload, error) {
	return s.fileRepo.GetPublicFileByID(ctx, id)
}

func (s *fileService) DeleteFile(ctx context.Context, id string) error {
	return s.fileRepo.DeleteFile(ctx, id)
}

func isAllowedFileType(ext string) bool {
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return allowedTypes[ext]
}
