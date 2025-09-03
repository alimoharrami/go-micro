package service

import (
	"context"
	"log"
	"notification/external/protos/userpb"
	"notification/internal/domain"
	"notification/internal/repository"
	"strconv"
)

type ChannelService struct {
	channelRep *repository.ChannelRepository
	userClient userpb.UserServiceClient
}

type CreateChannelInput struct {
	Name        string
	Description string
}

func NewChannelService(channelRep *repository.ChannelRepository, userClient userpb.UserServiceClient) *ChannelService {
	return &ChannelService{
		channelRep: channelRep,
		userClient: userClient,
	}
}

func (s *ChannelService) Create(ctx context.Context, input CreateChannelInput) (*domain.Channel, error) {
	channel := &domain.Channel{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := s.channelRep.Create(ctx, channel); err != nil {
		return nil, err
	}

	return channel, nil
}

func (s *ChannelService) Subscribe(ctx context.Context, userID int, channel string) error {
	return s.channelRep.Subscribe(ctx, userID, channel)
}

func (s *ChannelService) SendNotif(ctx context.Context, channle string, message string) error {
	subscriberIDs, err := s.channelRep.GetSubscriberIDs(ctx, channle)

	if err != nil {
		log.Printf("error send notif1 %v", err)
		return err
	}

	ids := make([]string, len(subscriberIDs))
	for i, id := range subscriberIDs {
		ids[i] = strconv.FormatUint(uint64(id), 10) // or similar conversion
	}

	resp, err := s.userClient.GetUsersByIDs(ctx, &userpb.GetUsersRequest{
		Ids: ids,
	})

	if err != nil {
		log.Printf("error send notif2 %v", err)
		return err
	}

	for _, user := range resp.Users {
		println("sending email to %v", user.Email)
	}

	return nil
}
