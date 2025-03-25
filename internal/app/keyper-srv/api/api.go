package api

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/service"
	"github.com/RomanAgaltsev/keyper/internal/config"
	"github.com/RomanAgaltsev/keyper/internal/logger/sl"
	"github.com/RomanAgaltsev/keyper/internal/model"
	"github.com/RomanAgaltsev/keyper/internal/pkg/auth"
	pb "github.com/RomanAgaltsev/keyper/pkg/keyper/v1"
	"github.com/RomanAgaltsev/keyper/pkg/transform"
)

const (
	msgInternalError      = "Internal error"
	msgLoginAlreadyTaken  = "Login has already been taken"
	msgWrongLoginPassword = "Wrong login/password"
	msgMissingUserID      = "Missing user ID"
)

var (
	_ UserService   = (*service.UserService)(nil)
	_ SecretService = (*service.SecretService)(nil)
)

type UserService interface {
	Register(ctx context.Context, user *model.User) error
	Login(ctx context.Context, user *model.User) error
}

type SecretService interface {
	Create(ctx context.Context, secret *model.Secret) (uuid.UUID, error)
	Update(ctx context.Context, userID uuid.UUID, secret *model.Secret) error
	Get(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) (*model.Secret, error)
	List(ctx context.Context, userID uuid.UUID) (model.Secrets, error)
	Delete(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) error
}

func NewUserAPI(log *slog.Logger, cfg *config.AppConfig, user UserService) pb.UserServiceServer {
	return &userAPI{
		log:  log,
		cfg:  cfg,
		user: user,
	}
}

type userAPI struct {
	log *slog.Logger
	cfg *config.AppConfig

	user UserService

	pb.UnimplementedUserServiceServer
}

func (a *userAPI) RegisterUserV1(ctx context.Context, request *pb.RegisterUserV1Request) (*pb.RegisterUserV1Response, error) {
	// TODO: observability

	const op = "userAPI.RegisterUser"

	user := transform.PbToUser(request.Credentials)

	err := a.user.Register(ctx, user)
	if err != nil && !errors.Is(err, service.ErrLoginTaken) {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, msgInternalError)
	}

	if errors.Is(err, service.ErrLoginTaken) {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.AlreadyExists, msgLoginAlreadyTaken)
	}

	// Generate JWT token
	ja := auth.NewAuth(a.cfg.SecretKey)
	_, tokenString, err := auth.NewJWTToken(ja, user, a.cfg.TokenTTL)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, msgInternalError)
	}

	response := pb.RegisterUserV1Response{
		Result: &pb.RegisterLoginResult{
			Token: tokenString,
		},
	}

	return &response, nil
}

func (a *userAPI) LoginUserV1(ctx context.Context, request *pb.LoginUserV1Request) (*pb.LoginUserV1Response, error) {
	const op = "userAPI.LoginUser"

	user := transform.PbToUser(request.Credentials)

	err := a.user.Login(ctx, user)
	if err != nil && !errors.Is(err, service.ErrLoginWrong) {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, msgInternalError)
	}

	if errors.Is(err, service.ErrLoginWrong) {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, msgWrongLoginPassword)
	}

	// Generate JWT token
	ja := auth.NewAuth(a.cfg.SecretKey)
	_, tokenString, err := auth.NewJWTToken(ja, user, a.cfg.TokenTTL)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, msgInternalError)
	}

	response := pb.LoginUserV1Response{
		Result: &pb.RegisterLoginResult{
			Token: tokenString,
		},
	}

	return &response, nil
}

func NewSecretAPI(log *slog.Logger, cfg *config.AppConfig, secret SecretService) pb.SecretServiceServer {
	return &secretAPI{
		log:    log,
		cfg:    cfg,
		secret: secret,
	}
}

type secretAPI struct {
	log *slog.Logger
	cfg *config.AppConfig

	secret SecretService

	pb.UnimplementedSecretServiceServer
}

func (a *secretAPI) CreateSecretV1(ctx context.Context, request *pb.CreateSecretV1Request) (*pb.CreateSecretV1Response, error) {
	const op = "secretAPI.CreateSecret"

	userID, err := auth.GetUserUID(ctx)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Unauthenticated, msgMissingUserID)
	}

	secret := transform.PbToSecret(request.Secret)
	secret.UserID = userID

	// TODO: add conflict handling
	secretID, err := a.secret.Create(ctx, secret)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, msgInternalError)
	}

	secretIDPb := secretID.String()
	errorPb := ""

	// TODO: return ID and error
	response := pb.CreateSecretV1Response{
		Result: &pb.CreateSecretV1Response_CreateSecretResult{
			Id:    &secretIDPb,
			Error: &errorPb,
		},
	}

	return &response, nil
}

func (a *secretAPI) UpdateSecretV1(ctx context.Context, request *pb.UpdateSecretV1Request) (*pb.UpdateSecretV1Response, error) {
	const op = "secretAPI.UpdateSecret"

	userID, err := auth.GetUserUID(ctx)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Unauthenticated, msgMissingUserID)
	}

	secret := transform.PbToSecret(request.Secret)

	err = a.secret.Update(ctx, userID, secret)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, msgInternalError)
	}

	// TODO: transform error
	errorPb := ""

	response := pb.UpdateSecretV1Response{
		Result: &pb.UpdateSecretResult{
			Error: &errorPb,
		},
	}

	return &response, nil
}

func (a *secretAPI) UpdateSecretsDataV1(stream grpc.ClientStreamingServer[pb.UpdateSecretsDataV1Request, pb.UpdateSecretsDataV1Response]) error {
	// TODO: add secrets ownership check
	return nil
}

func (a *secretAPI) GetSecretV1(ctx context.Context, request *pb.GetSecretV1Request) (*pb.GetSecretV1Response, error) {
	const op = "secretAPI.GetSecret"

	userID, err := auth.GetUserUID(ctx)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Unauthenticated, msgMissingUserID)
	}

	secretID, err := uuid.Parse(request.Id)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, msgInternalError)
	}

	secret, err := a.secret.Get(ctx, userID, secretID)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, msgInternalError)
	}

	secretPb := transform.SecretToPb(secret)
	errorPb := ""

	// TODO: return secret and error
	response := pb.GetSecretV1Response{
		Result: &pb.GetSecretV1Response_GetSecretResult{
			Secret: secretPb,
			Error:  &errorPb,
		},
	}

	return &response, nil
}

func (a *secretAPI) GetSecretsDataV1(request *pb.GetSecretsDataV1Request, stream grpc.ServerStreamingServer[pb.GetSecretsDataV1Response]) error {
	// TODO: add secrets ownership check
	return nil
}

func (a *secretAPI) ListSecretsV1(ctx context.Context, _ *emptypb.Empty) (*pb.ListSecretsV1Response, error) {
	const op = "secretAPI.ListSecrets"

	userID, err := auth.GetUserUID(ctx)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Unauthenticated, msgMissingUserID)
	}

	_, err = a.secret.List(ctx, userID)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, msgInternalError)
	}

	// TODO: transform list of secrets and error
	// TODO: return list of secrets and error

	return nil, nil
}

func (a *secretAPI) DeleteSecretV1(ctx context.Context, request *pb.DeleteSecretV1Request) (*pb.DeleteSecretV1Response, error) {
	const op = "secretAPI.DeleteSecret"

	userID, err := auth.GetUserUID(ctx)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Unauthenticated, msgMissingUserID)
	}

	secretID, err := uuid.Parse(request.Id)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, msgInternalError)
	}

	// TODO: add secrets ownership check
	// TODO: add errors messages
	err = a.secret.Delete(ctx, userID, secretID)
	if err != nil {
		a.log.Error(op, sl.Err(err))
		return nil, status.Error(codes.Internal, msgInternalError)
	}

	// TODO: transform error
	errorPb := ""

	// TODO: return error
	response := pb.DeleteSecretV1Response{
		Result: &pb.DeleteSecretV1Response_DeleteSecretResult{
			Error: &errorPb,
		},
	}

	return &response, nil
}

func Register(gRPCServer *grpc.Server, userAPI pb.UserServiceServer, secretAPI pb.SecretServiceServer) {
	pb.RegisterUserServiceServer(gRPCServer, userAPI)
	pb.RegisterSecretServiceServer(gRPCServer, secretAPI)
}
