package api

import (
	"context"
	"log/slog"

	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/RomanAgaltsev/keyper/internal/logger/sl"
	"github.com/RomanAgaltsev/keyper/internal/model"
	pb "github.com/RomanAgaltsev/keyper/pkg/keyper/v1"
)

type UserService interface {
	Register(ctx context.Context, user model.User) error
	Login(ctx context.Context, user model.User) error
}

type SecretService interface {
	Create(ctx context.Context, secret model.Secret) (uuid.UUID, error)
	Get(ctx context.Context, secretID uuid.UUID) (model.Secret, error)
	List(ctx context.Context, user model.User) (model.Secrets, error)
	Update(ctx context.Context, secret model.Secret) error
	Delete(ctx context.Context, secretID uuid.UUID) error
}

func NewUserAPI(log *slog.Logger, user UserService) pb.UserServiceServer {
	return &userAPI{
		log:  log,
		user: user,
	}
}

type userAPI struct {
	log  *slog.Logger
	user UserService

	pb.UnimplementedUserServiceServer
}

func (a *userAPI) RegisterUserV1(ctx context.Context, request *pb.RegisterUserV1Request) (*pb.RegisterUserV1Response, error) {
	// TODO: move validation to middleware
	if err := protovalidate.Validate(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	const op = "userAPI.RegisterUser"

	// TODO: transform user from request
	user := model.User{}

	// TODO: add conflict handling
	// TODO: add errors messages
	err := a.user.Register(ctx, user)
	if err != nil {
		a.log.Error(op, sl.Err(err))
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

	const op = "userAPI.LoginUser"

	// TODO: transform user from request
	user := model.User{}

	// TODO: add errors messages
	err := a.user.Login(ctx, user)
	if err != nil {
		a.log.Error(op, sl.Err(err))
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

func NewSecretAPI(log *slog.Logger, secret SecretService) pb.SecretServiceServer {
	return &secretAPI{
		log:    log,
		secret: secret,
	}
}

type secretAPI struct {
	log    *slog.Logger
	secret SecretService

	pb.UnimplementedSecretServiceServer
}

func (a *secretAPI) CreateSecretV1(ctx context.Context, request *pb.CreateSecretV1Request) (*pb.CreateSecretV1Response, error) {
	// TODO: move validation to middleware
	if err := protovalidate.Validate(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	const op = "secretAPI.CreateSecret"

	// TODO: transform secret from request
	secret := model.Secret{}

	// TODO: add conflict handling
	// TODO: add errors messages
	secretID, err := a.secret.Create(ctx, secret)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, "please look at logs")
	}

	// TODO: transform ID and error

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

	const op = "secretAPI.GetSecret"

	// TODO: transform secret ID from request
	secretID := uuid.New()

	// TODO: add errors messages
	secret, err := a.secret.Get(ctx, secretID)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, "please look at logs")
	}

	// TODO: transform secret and error

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

	const op = "secretAPI.ListSecrets"

	// TODO: transform user from request
	user := model.User{}

	// TODO: add errors messages
	secrets, err := a.secret.List(stream.Context(), user)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return status.Error(codes.Internal, "please look at logs")
	}

	// TODO: transform list of secrets and error
	// TODO: return list of secrets and error

	return nil
}

func (a *secretAPI) UpdateSecretV1(ctx context.Context, request *pb.UpdateSecretV1Request) (*pb.UpdateSecretV1Response, error) {
	// TODO: move validation to middleware
	if err := protovalidate.Validate(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	const op = "secretAPI.UpdateSecret"

	// TODO: transform secret from request
	secret := model.Secret{}

	// TODO: add errors messages
	err := a.secret.Update(ctx, secret)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, "please look at logs")
	}

	// TODO: transform error

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

	const op = "secretAPI.DeleteSecret"

	// TODO: transform secret ID from request
	secretID := uuid.New()

	// TODO: add errors messages
	err := a.secret.Delete(ctx, secretID)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, "please look at logs")
	}

	// TODO: transform error

	// TODO: return error
	response := pb.DeleteSecretV1Response{
		Result: &pb.DeleteSecretV1Response_DeleteSecretResult{
			Error: nil,
		},
	}

	return &response, nil
}

func Register(gRPCServer *grpc.Server, userAPI pb.UserServiceServer, secretAPI pb.SecretServiceServer) {
	pb.RegisterUserServiceServer(gRPCServer, userAPI)
	pb.RegisterSecretServiceServer(gRPCServer, secretAPI)
}
