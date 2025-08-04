package domain

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null"` // e.g. "admin"
	Slug        string `gorm:"uniqueIndex;not null"` // e.g. "admin"
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
