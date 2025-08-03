package migrations

import (
	"log"
	"user/internal/domain"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
}
