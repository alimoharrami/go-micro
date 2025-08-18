package usersr

import (
	"context"
	"log"
	"strconv"
	"user/external/protos/userpb"
	"user/internal/service"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserServer(userService *service.UserService) userpb.UserServiceServer {
	return &UserServer{userService: userService}
}

func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	idUint, _ := strconv.ParseUint(req.Id, 10, 64)
	user, err := s.userService.GetByID(ctx, uint(idUint))

	if err != nil {
		log.Fatalf("error fetching data: %v", err)
	}

	return &userpb.GetUserResponse{
		Id:   req.Id,
		Name: user.FirstName,
	}, nil
}
