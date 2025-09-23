package seeders

import (
	"context"
	"fmt"
	"math/rand"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/pkg/database"
	"time"

	"github.com/segmentio/ksuid"
)

func SeedFileUploads(ctx context.Context) error {
	db := database.GetDB()

	// Get existing users
	var users []domain.User
	if err := db.NewSelect().Model(&users).Scan(ctx); err != nil {
		return err
	}
	if len(users) == 0 {
		return fmt.Errorf("no users found, please seed users first")
	}

	contentTypes := []string{
		"image/jpeg",
		"image/png",
		"application/pdf",
		"text/plain",
		"application/msword",
	}

	modules := []string{
		"blog",
		"profile",
		"document",
		"gallery",
	}

	// Create 10 file uploads
	for i := 1; i <= 10; i++ {
		fileUpload := &domain.FileUpload{
			ID:          ksuid.New().String(),
			FileName:    fmt.Sprintf("sample-file-%d.%s", i, getExtension(contentTypes[rand.Intn(len(contentTypes))])),
			FilePath:    fmt.Sprintf("/uploads/%d/%s", time.Now().Year(), fmt.Sprintf("file-%d", i)),
			FileSize:    int64(rand.Intn(10000000)), // Random size up to 10MB
			ContentType: contentTypes[rand.Intn(len(contentTypes))],
			Module:     modules[rand.Intn(len(modules))],
			AuthorID:   users[rand.Intn(len(users))].ID,
			IsPublic:   rand.Float32() < 0.7, // 70% chance of being public
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if _, err := db.NewInsert().Model(fileUpload).Exec(ctx); err != nil {
			return err
		}
	}

	fmt.Println("✓ File uploads seeded")
	return nil
}

func getExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "application/pdf":
		return "pdf"
	case "text/plain":
		return "txt"
	case "application/msword":
		return "doc"
	default:
		return "bin"
	}
}