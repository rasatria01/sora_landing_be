package seeders

import (
	"context"
	"fmt"
	"math/rand"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/utils"

	"github.com/segmentio/ksuid"
)

func SeedTags(ctx context.Context) error {
	db := database.GetDB()
	tags := []string{
		"Golang",
		"JavaScript",
		"Python",
		"React",
		"Vue.js",
		"Node.js",
		"Docker",
		"Kubernetes",
		"AWS",
		"Git",
	}
	var users []domain.User
	if err := db.NewSelect().Model(&users).Scan(ctx); err != nil {
		return err
	}
	for _, name := range tags {
		tag := &domain.Tag{
			BaseEntity: domain.BaseEntity{
				ID: ksuid.New().String(),
			},
			Name:        name,
			Slug:        utils.Slugify(name),
			CreatedByID: users[rand.Intn(len(users))].ID,
			EditedByID:  users[rand.Intn(len(users))].ID,
		}

		if _, err := db.NewInsert().Model(tag).Exec(ctx); err != nil {
			return err
		}
	}

	fmt.Println("✓ Tags seeded")
	return nil
}
