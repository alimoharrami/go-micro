package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
}

func NewServer(register func(*grpc.Server)) *GRPCServer {
	s := grpc.NewServer()
	register(s)
	return &GRPCServer{server: s}
}

func (g *GRPCServer) Start(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("gRPC server running on %s", address)
	if err := g.server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
