package api

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (a *userAPI) RegisterUserV1(ctx context.Context, request *pb.RegisterUserV1Request) (*pb.RegisterUserV1Response, error) {
	// TODO: move validation to middleware
	if err := protovalidate.Validate(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: add conflict handling
	// TODO: add errors messages
	err := a.user.Register(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "please look at logs")
	}

	// TODO: return token
	response := pb.RegisterUserV1Response{
		Result: &pb.RegisterLoginResult{
			Token: "",
		},
	}

	return &response, nil
}

func (a *userAPI) LoginUserV1(ctx context.Context, request *pb.LoginUserV1Request) (*pb.LoginUserV1Response, error) {
	// TODO: move validation to middleware
	if err := protovalidate.Validate(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: add errors messages
	err := a.user.Login(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "please look at logs")
	}

	// TODO: return token
	response := pb.LoginUserV1Response{
		Result: &pb.RegisterLoginResult{
			Token: "",
		},
	}

	return &response, nil
}

type secretAPI struct {
	pb.UnimplementedSecretServiceServer
	secret SecretService
}

func (a *secretAPI) CreateSecretV1(ctx context.Context, request *pb.CreateSecretV1Request) (*pb.CreateSecretV1Response, error) {
	// TODO: move validation to middleware
	if err := protovalidate.Validate(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: add conflict handling
	// TODO: add errors messages
	err := a.secret.Create(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "please look at logs")
	}

	// TODO: return ID and error
	response := pb.CreateSecretV1Response{
		Result: &pb.CreateSecretV1Response_CreateSecretResult{
			Id:    nil,
			Error: nil,
		},
	}

	return &response, nil
}

func (a *secretAPI) GetSecretV1(ctx context.Context, request *pb.GetSecretV1Request) (*pb.GetSecretV1Response, error) {
	// TODO: move validation to middleware
	if err := protovalidate.Validate(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: add errors messages
	err := a.secret.Get(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "please look at logs")
	}

	// TODO: return secret and error
	response := pb.GetSecretV1Response{
		Result: &pb.GetSecretV1Response_GetSecretResult{
			Secret: nil,
			Error:  nil,
		},
	}

	return &response, nil
}

func (a *secretAPI) ListSecretsV1(request *pb.ListSecretsV1Request, stream grpc.ServerStreamingServer[pb.ListSecretsV1Response]) error {
	// TODO: move validation to middleware
	if err := protovalidate.Validate(request); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: add errors messages
	err := a.secret.List(context.Background())
	if err != nil {
		return status.Error(codes.Internal, "please look at logs")
	}

	// TODO: return list of secrets and error

	return nil
}

func (a *secretAPI) UpdateSecretV1(ctx context.Context, request *pb.UpdateSecretV1Request) (*pb.UpdateSecretV1Response, error) {
	// TODO: move validation to middleware
	if err := protovalidate.Validate(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: add errors messages
	err := a.secret.Update(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "please look at logs")
	}

	// TODO: return error
	response := pb.UpdateSecretV1Response{
		Result: &pb.UpdateSecretV1Response_UpdateSecretResult{
			Error: nil,
		},
	}

	return &response, nil
}

func (a *secretAPI) DeleteSecretV1(ctx context.Context, request *pb.DeleteSecretV1Request) (*pb.DeleteSecretV1Response, error) {
	// TODO: move validation to middleware
	if err := protovalidate.Validate(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: add errors messages
	err := a.secret.Delete(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "please look at logs")
	}

	// TODO: return error
	response := pb.DeleteSecretV1Response{
		Result: &pb.DeleteSecretV1Response_DeleteSecretResult{
			Error: nil,
		},
	}

	return &response, nil
}

func Register(gRPCServer *grpc.Server, user UserService, secret SecretService) {
	pb.RegisterUserServiceServer(gRPCServer, &userAPI{user: user})
	pb.RegisterSecretServiceServer(gRPCServer, &secretAPI{secret: secret})
}
