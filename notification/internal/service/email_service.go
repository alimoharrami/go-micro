package service

import (
	"log"
	"notification/external/protos/userpb"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmailService struct {
	userClient userpb.UserServiceClient
}

func NewEmailService(userClient userpb.UserServiceClient) *EmailService {
	return &EmailService{userClient}
}

func (s *EmailService) SendEmailUserID(c *gin.Context, UserID int, message string) {
	user, err := s.userClient.GetUser(c, &userpb.GetUserRequest{Id: strconv.Itoa(UserID)})
	if err != nil {
		log.Println("error in notificaiotn grpc %v", err)
	}
	log.Println(message + "" + user.Name)
}
