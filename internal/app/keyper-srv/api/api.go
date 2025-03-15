package api

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/RomanAgaltsev/keyper/pkg/keyper/v1"
)

type UserService interface {
	Register(ctx context.Context) error
	Login(ctx context.Context) error
}

type SecretService interface {
	Create(ctx context.Context) error
	Get(ctx context.Context) error
	List(ctx context.Context) error
	Update(ctx context.Context) error
	Delete(ctx context.Context) error
}

type userAPI struct {
	pb.UnimplementedUserServiceServer
	user UserService
}

type secretAPI struct {
	pb.UnimplementedSecretServiceServer
	secret SecretService
}

func Register(gRPCServer *grpc.Server, user UserService, secret SecretService) {
	pb.RegisterUserServiceServer(gRPCServer, &userAPI{user: user})
	pb.RegisterSecretServiceServer(gRPCServer, &secretAPI{secret: secret})
}
