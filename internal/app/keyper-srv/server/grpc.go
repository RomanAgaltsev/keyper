package server

import (
	"log/slog"

	"google.golang.org/grpc"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/api"
)

func NewGRPCServer(log *slog.Logger, userService api.UserService, secretService api.SecretService) *grpc.Server {
	// gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor( recovery.UnaryServerInterceptor(recoveryOpts...), logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),))
	server := grpc.NewServer()

	userAPI := api.NewUserAPI(log, userService)
	secretAPI := api.NewSecretAPI(log, secretService)

	api.Register(server, userAPI, secretAPI)

	return server
}
