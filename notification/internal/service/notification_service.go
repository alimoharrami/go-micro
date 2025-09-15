package service

import (
	"context"
	"fmt"
	"notification/internal/domain"
	"notification/internal/hub"
	"notification/internal/repository"
	"strconv"
)

type CreateNotificationInput struct {
	Title       string
	Description string
	userID      uint
}

type NotificaitonSerivce struct {
	notifRep  *repository.NotificationRepository
	notifHelp *hub.NotificationHelper
}

func NewNotificaitonSerivce(
	notifRep *repository.NotificationRepository,
	notifHelp *hub.NotificationHelper,
) *NotificaitonSerivce {
	return &NotificaitonSerivce{
		notifRep:  notifRep,
		notifHelp: notifHelp,
	}
}

func (s *NotificaitonSerivce) Store(ctx context.Context, input CreateNotificationInput) (*domain.Notification, error) {
	notification := &domain.Notification{
		Title:       input.Title,
		Description: input.Description,
		UserID:      input.userID,
	}

	if err := s.notifRep.Store(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *NotificaitonSerivce) GetUserNotifCount(ctx context.Context, userID uint) (int64, error) {
	return s.notifRep.GetUserNotifCount(ctx, int(userID))
}

func (s *NotificaitonSerivce) BroadcastCount(ctx context.Context, userID int) {
	count, err := s.notifRep.GetUserNotifCount(ctx, int(userID))
	if err != nil {
		fmt.Printf("error in broadcasting count %v", err)
	}
	countString := strconv.FormatInt(count, 10)
	s.notifHelp.Broadcast(countString)
}
