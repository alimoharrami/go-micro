package migrations

import (
	"log"
	"notification/internal/domain"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&domain.Channel{},
		&domain.ChannelUser{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
}
