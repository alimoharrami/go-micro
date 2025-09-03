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
		User: &userpb.User{
			Id:    req.Id, // assuming user.ID is string, else convert accordingly
			Name:  user.FirstName,
			Email: user.Email, // assuming you want to include email
		},
	}, nil
}

func (s *UserServer) GetUsersByIDs(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	var uintIDs []uint
	for _, idStr := range req.Ids {
		idUint64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return nil, err // handle parsing error properly
		}
		uintIDs = append(uintIDs, uint(idUint64)) // convert uint64 to uint here
	}

	users, err := s.userService.GetUserListByIDs(ctx, uintIDs)
	if err != nil {
		return nil, err
	}

	// Convert your user domain model to protobuf User messages
	var pbUsers []*userpb.User
	for _, u := range users {
		pbUser := &userpb.User{
			Id:    strconv.FormatUint(uint64(u.ID), 10),
			Name:  u.FirstName,
			Email: u.Email,
		}
		pbUsers = append(pbUsers, pbUser)
	}

	// Return the response with the slice of users
	return &userpb.GetUsersResponse{
		Users: pbUsers,
	}, nil
}
