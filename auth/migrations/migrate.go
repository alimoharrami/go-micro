package migrations

import (
	"auth/internal/domain"
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
		&domain.Permission{},
		&domain.RolePermission{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
}
