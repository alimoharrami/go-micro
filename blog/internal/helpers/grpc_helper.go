package helpers

import (
	"go-blog/external/protos/userpb"

	"github.com/alimoharrami/go-micro/pkg/grpc"
)

func InitGRPC() userpb.UserServiceClient {
	conn := grpc.NewClientConn("user-service:50051")
	// Note: We are not closing the connection here. 
	// In a real production app, we should use fx.Hook to close it on shutdown.

	client := userpb.NewUserServiceClient(conn)
	return client
}
