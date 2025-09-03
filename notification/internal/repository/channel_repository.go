package repository

import (
	"context"
	"errors"
	"notification/internal/domain"

	"gorm.io/gorm"
)

type ChannelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) Create(ctx context.Context, channel *domain.Channel) error {
	return r.db.WithContext(ctx).Create(channel).Error
}

func (r *ChannelRepository) Subscribe(ctx context.Context, userID int, channelName string) error {
	var channel domain.Channel

	if err := r.db.Where("name = ?", channelName).First(&channel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("channel not found")
		}
		return err
	}

	var existing domain.ChannelUser
	err := r.db.Where("channel_id = ? AND user_id = ?", channel.ID, userID).First(&existing).Error
	if err == nil {
		return errors.New("user already subscribed")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	sub := domain.ChannelUser{
		ChannelID: channel.ID,
		UserID:    uint(userID),
	}
	if err := r.db.Create(&sub).Error; err != nil {
		return err
	}
	return nil
}

func (r *ChannelRepository) GetSubscriberIDs(ctx context.Context, channelName string) ([]uint, error) {
	var channel domain.Channel
	err := r.db.Where("name = ?", channelName).First(&channel).Error
	if err != nil {
		return nil, errors.New("channel not exitsts")
	}

	channelID := channel.ID

	var subscribers []domain.ChannelUser

	errr := r.db.Where("channel_id = ?", channelID).Find(&subscribers).Error

	if errr != nil {
		return nil, err
	}

	subscriberIDs := make([]uint, len(subscribers))

	for i, s := range subscribers {
		subscriberIDs[i] = s.UserID
	}

	return subscriberIDs, nil
}
