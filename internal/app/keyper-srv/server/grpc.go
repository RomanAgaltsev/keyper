package server

import (
	"google.golang.org/grpc"
	"log/slog"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/api"
)

func NewGRPCServer(log *slog.Logger, user api.UserService, secret api.SecretService) *grpc.Server {
	//gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor( recovery.UnaryServerInterceptor(recoveryOpts...), logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),))
	server := grpc.NewServer()

	api.Register(server, user, secret)

	return server
}
