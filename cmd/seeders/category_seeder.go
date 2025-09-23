package seeders

import (
	"context"
	"fmt"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/utils"

	"github.com/segmentio/ksuid"
)

func SeedCategories(ctx context.Context) error {
	db := database.GetDB()
	categories := []string{
		"Technology",
		"Programming",
		"Web Development",
		"Mobile Development",
		"DevOps",
		"Cloud Computing",
		"Artificial Intelligence",
		"Machine Learning",
		"Cybersecurity",
		"Data Science",
	}

	for _, name := range categories {
		category := &domain.Category{
			BaseEntity: domain.BaseEntity{
				ID: ksuid.New().String(),
			},
			Name: name,
			Slug: utils.Slugify(name),
		}

		if _, err := db.NewInsert().Model(category).Exec(ctx); err != nil {
			return err
		}
	}

	fmt.Println("✓ Categories seeded")
	return nil
}
