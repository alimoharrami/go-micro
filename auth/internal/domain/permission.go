package domain

import (
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Key         string `gorm:"uniqueIndex;not null"` // e.g. "read:users", "create:accounts"
	Description string
}
