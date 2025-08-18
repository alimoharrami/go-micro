package helpers

import (
	"context"
	"go-blog/external/protos/userpb"
	"log"
	"time"

	"github.com/alimoharrami/go-micro/pkg/grpc"
)

func InitGRPC() {
	conn := grpc.NewClientConn("user-service:50051")
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &userpb.GetUserRequest{Id: "1"})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("User: %s - %s", resp.Id, resp.Name)
}
