package seeders

import (
	"context"
	"fmt"
	"sora_landing_be/cmd/constants"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/pkg/authentication"
	"sora_landing_be/pkg/database"
)

func SeedUsers(ctx context.Context) error {
	db := database.GetDB()

	// Create 10 users with authentication
	for i := 1; i <= 10; i++ {

		// Create password hash
		hashedPassword, err := authentication.HashPassword(fmt.Sprintf("password%d", i))
		if err != nil {
			return err
		}

		user := &domain.User{
			Name:   fmt.Sprintf("User %d", i),
			Email:  fmt.Sprintf("user%d@example.com", i),
			Status: "Active",
			Roles:  []constants.UserRole{constants.UserRoleUser},
		}

		// First user is admin
		if i == 1 {
			user.Name = "Admin User"
			user.Email = "admin@example.com"
			user.Roles = []constants.UserRole{constants.UserRoleAdmin}
		}

		auth := &domain.Authentication{
			UserID:   user.ID,
			Password: hashedPassword,
		}

		if _, err := db.NewInsert().Model(user).Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewInsert().Model(auth).Exec(ctx); err != nil {
			return err
		}
	}

	fmt.Println("✓ Users seeded")
	return nil
}
