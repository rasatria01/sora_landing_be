package seeders

import (
	"context"
	"database/sql"
	"log"
	"sora_landing_be/pkg/database"

	"github.com/uptrace/bun"
)

// RunSeeder runs all seeders in the correct order with a database transaction
func RunSeeder(ctx context.Context) error {
	db := database.GetDB()
	return db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		log.Println("Starting database seeding...")

		// Order matters due to foreign key constraints
		seeders := []func(context.Context) error{
			SeedUsers,        // First, as other tables reference users
			SeedCategories,   // Second, as blogs need categories
			SeedTags,         // Third, as blogs need tags
			SeedBlogArtikels, // Fourth, depends on users, categories, and tags
			SeedFileUploads,  // Last, depends on users
		}

		for _, seeder := range seeders {
			if err := seeder(ctx); err != nil {
				return err
			}
		}

		log.Println("Database seeding completed successfully")
		return nil
	})
}
