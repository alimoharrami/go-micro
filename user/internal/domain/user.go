package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Active    bool   `gorm:"default:true"`
	LastLogin time.Time
	Role      string `gorm:"size:50;default:user"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
