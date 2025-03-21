package server

import (
	"log/slog"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	logging_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	recovery_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/api"
	"github.com/RomanAgaltsev/keyper/internal/config"
	"github.com/RomanAgaltsev/keyper/internal/logger/sl"
)

func NewGRPCServer(log *slog.Logger, cfg *config.GRPCConfig, userService api.UserService, secretService api.SecretService) *grpc.Server {
	const op = "server.NewGRPCServer"

	loggerOpts := InterceptorLoggerOpts()
	recoveryOpts := InterceptorRecoveryOpts()

	validator, err := protovalidate.New()
	if err != nil {
		log.Error(op, sl.Err(err))
	}

	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: cfg.MaxConnectionIdle,
			Timeout:           cfg.Timeout,
			MaxConnectionAge:  cfg.MaxConnectionAge,
			Time:              cfg.Timeout,
		}),
		grpc.ChainUnaryInterceptor(
			logging_middleware.UnaryServerInterceptor(InterceptorLogger(log), loggerOpts...),
			recovery_middleware.UnaryServerInterceptor(recoveryOpts...),
			protovalidate_middleware.UnaryServerInterceptor(validator),

		),
		grpc.ChainStreamInterceptor(
			logging_middleware.StreamServerInterceptor(InterceptorLogger(log), loggerOpts...),
			recovery_middleware.StreamServerInterceptor(recoveryOpts...),
			protovalidate_middleware.StreamServerInterceptor(validator),

		),
	)

	userAPI := api.NewUserAPI(log, userService)
	secretAPI := api.NewSecretAPI(log, secretService)

	api.Register(server, userAPI, secretAPI)

	return server
}
