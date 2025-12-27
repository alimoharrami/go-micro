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
		&domain.Notification{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// Seed default channel
	var count int64
	db.Model(&domain.Channel{}).Where("name = ?", "main").Count(&count)
	if count == 0 {
		mainChannel := domain.Channel{
			Name:        "main",
			Description: "Main broadcast channel",
		}
		if err := db.Create(&mainChannel).Error; err != nil {
			log.Printf("Failed to seed main channel: %v", err)
		} else {
			log.Println("Successfully seeded main channel")
		}
	}
}
