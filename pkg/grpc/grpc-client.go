package grpc

import (
	"log"

	"google.golang.org/grpc"
)

func NewClientConn(address string) *grpc.ClientConn {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	return conn
}
