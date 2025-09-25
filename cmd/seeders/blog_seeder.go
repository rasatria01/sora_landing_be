package seeders

import (
	"context"
	"fmt"
	"math/rand"
	"sora_landing_be/cmd/constants"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/utils"
	"time"
)

func SeedBlogArtikels(ctx context.Context) error {
	db := database.GetDB()

	// Get existing users, categories, and tags
	var users []domain.User
	var categories []domain.Category
	var tags []domain.Tag

	if err := db.NewSelect().Model(&users).Scan(ctx); err != nil {
		return err
	}
	if err := db.NewSelect().Model(&categories).Scan(ctx); err != nil {
		return err
	}
	if err := db.NewSelect().Model(&tags).Scan(ctx); err != nil {
		return err
	}

	// Check if we have users, categories and tags
	if len(users) == 0 {
		return fmt.Errorf("no users found, please seed users first")
	}
	if len(categories) == 0 {
		return fmt.Errorf("no categories found, please seed categories first")
	}
	if len(tags) == 0 {
		return fmt.Errorf("no tags found, please seed tags first")
	}

	// Create 10 blog articles
	for i := 1; i <= 10; i++ {
		title := fmt.Sprintf("Sample Blog asda hahah Article %d", i)

		article := &domain.BlogArtikel{

			Title:       title,
			Slug:        utils.Slugify(title),
			Content:     fmt.Sprintf("This is the content for article %d. It contains detailed information about the topic.", i),
			Excerpt:     fmt.Sprintf("Brief excerpt for article %d", i),
			ImageURL:    fmt.Sprintf("https://example.com/images/article-%d.jpg", i),
			CategoryID:  categories[rand.Intn(len(categories))].ID,
			AuthorID:    users[rand.Intn(len(users))].ID,
			Status:      constants.StatusPublished,
			Views:       int64(rand.Intn(1000)),
			PublishedAt: time.Now(),
		}

		if _, err := db.NewInsert().Model(article).Returning("id").Exec(ctx); err != nil {
			return err
		}

		// Add 2-4 random tags for each article
		numTags := rand.Intn(3) + 1
		selectedTags := make(map[string]bool)

		for range numTags {
			tagID := tags[rand.Intn(len(tags))].ID
			if !selectedTags[tagID] {
				articleTag := &domain.ArticleTag{
					ArticleID: article.ID,
					TagID:     tagID,
				}
				if _, err := db.NewInsert().Model(articleTag).Exec(ctx); err != nil {
					return err
				}
				selectedTags[tagID] = true
			}
		}
	}

	fmt.Println("✓ Blog articles seeded")
	return nil
}
