package domain

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	Title       string
	Description string
	UserID      uint `json:"user_id"`
}
