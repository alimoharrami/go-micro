package service

import (
	"context"
	"log"
	"notification/external/protos/userpb"
	"strconv"
)

type EmailService struct {
	userClient userpb.UserServiceClient
}

func NewEmailService(userClient userpb.UserServiceClient) *EmailService {
	return &EmailService{userClient}
}

func (s *EmailService) SendEmailUserID(c context.Context, UserID int, message string) {
	user, err := s.userClient.GetUser(c, &userpb.GetUserRequest{Id: strconv.Itoa(UserID)})
	if err != nil {
		log.Printf("error in notificaiotn grpc %v", err)
	}
	log.Println(message + " " + user.Name)
}
