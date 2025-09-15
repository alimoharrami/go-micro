package repository

import (
	"context"
	"notification/internal/domain"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Store(ctx context.Context, notification *domain.Notification) error {

	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *NotificationRepository) GetUserNotifCount(ctx context.Context, userID int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	return count, err
}
