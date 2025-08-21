package domain

import "gorm.io/gorm"

type Channel struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null"` // e.g. "admin"
	Description string
}
