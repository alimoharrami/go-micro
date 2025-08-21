package domain

import "gorm.io/gorm"

type ChannelUser struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	ChannelID uint
	UserID    uint
}
