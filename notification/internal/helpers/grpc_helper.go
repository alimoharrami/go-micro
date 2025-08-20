package helpers

import (
	"notification/external/protos/userpb"

	"github.com/alimoharrami/go-micro/pkg/grpc"
)

func InitGRPC() userpb.UserServiceClient {
	conn := grpc.NewClientConn("user-service:50051")
	// defer conn.Close()

	client := userpb.NewUserServiceClient(conn)
	return client
}
