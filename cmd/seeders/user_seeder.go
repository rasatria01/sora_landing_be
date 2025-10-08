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

	// Check if superadmin exists
	exists, err := db.NewSelect().Model((*domain.User)(nil)).
		Where("roles @> ?", "{SuperAdmin}").
		Exists(ctx)
	if err != nil {
		return fmt.Errorf("error checking superadmin: %v", err)
	}

	if exists {
		fmt.Println("✓ Superadmin already exists")
		return nil
	}

	// Create superadmin
	hashedPassword, err := authentication.HashPassword("admin123") // You might want to change this password
	if err != nil {
		return err
	}

	user := &domain.User{
		Name:   "Admin User",
		Email:  "admin@example.com",
		Status: "Active",
		Roles:  []constants.UserRole{constants.UserRoleSuperAdmin},
	}

	// Insert user first
	if _, err := db.NewInsert().Model(user).Exec(ctx); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	// Create authentication record
	auth := &domain.Authentication{
		UserID:   user.ID,
		Password: hashedPassword,
	}

	// Insert authentication
	if _, err := db.NewInsert().Model(auth).Exec(ctx); err != nil {
		return fmt.Errorf("failed to create authentication: %v", err)
	}

	// Verify the user was created properly
	var auth2 domain.Authentication
	err = db.NewSelect().Model(&auth2).
		Where("user_id = ?", user.ID).
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("failed to verify authentication creation: %v", err)
	}

	fmt.Println("✓ Users seeded")
	return nil
}
