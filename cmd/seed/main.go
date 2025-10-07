package main

import (
	"context"
	"flag"
	"log"
	"sora_landing_be/cmd/seeders"
	"sora_landing_be/pkg/config"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/logger"
)

func main() {
	// Parse command line flags
	table := flag.String("table", "", "Specific table to seed (users, categories, tags, blogs, files)")
	flag.Parse()

	// Load environment variables first

	// Initialize configuration and database
	cfg := config.LoadConfig()
	logger.NewZapLogger(cfg.Logger)
	database.InitDB(cfg.Database)
	ctx := context.Background()

	// Get context
	// Run specific seeder or all
	var err error
	switch *table {
	case "users":
		err = seeders.SeedUsers(ctx)
	case "categories":
		err = seeders.SeedCategories(ctx)
	case "tags":
		err = seeders.SeedTags(ctx)
	case "blogs":
		err = seeders.SeedBlogArtikels(ctx)
	case "":
		err = seeders.RunSeeder(ctx)
	default:
		log.Fatalf("Unknown table: %s", *table)
	}

	if err != nil {
		log.Fatalf("Error seeding database: %v", err)
	}
}
