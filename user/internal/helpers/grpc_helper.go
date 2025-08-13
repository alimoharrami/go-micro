package helpers

import (
	"context"
	"user/external/protos/userpb"

	mygrpc "github.com/alimoharrami/go-micro/pkg/grpc"
	"google.golang.org/grpc"
)

type userServer struct {
	userpb.UnimplementedUserServiceServer
}

func (s *userServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	return &userpb.GetUserResponse{
		Id:   req.Id,
		Name: "Ali Moharrami",
	}, nil
}

func InitilaizeGRPC() {
	go func() {
		serverGrpc := mygrpc.NewServer(func(s *grpc.Server) {
			userpb.RegisterUserServiceServer(s, &userServer{})
		})

		serverGrpc.Start(":50051")
	}()
}
