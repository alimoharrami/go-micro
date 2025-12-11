package helpers

import (
	"context"
	"log"
	"net"
	"runtime/debug"
	"user/external/protos/userpb"
	"user/internal/repository"
	"user/internal/service"

	"user/external/grpc/usersr"

	"github.com/alimoharrami/go-micro/pkg/rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func InitilaizeGRPC(db *gorm.DB, publisher rabbitmq.IPublisher) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, publisher)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] gRPC server: %v\n%s", r, debug.Stack())
			}
		}()

		server := grpc.NewServer(
			grpc.UnaryInterceptor(RecoveryInterceptor),
		)
		userpb.RegisterUserServiceServer(server, usersr.NewUserServer(userService))

		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()
}

func RecoveryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[PANIC] in %s: %v\n", info.FullMethod, r)
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()

	// normal flow
	return handler(ctx, req)
}
