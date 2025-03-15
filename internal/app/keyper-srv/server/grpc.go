package server

import (
	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/api"
	"google.golang.org/grpc"

	"github.com/RomanAgaltsev/keyper/internal/config"
)

func NewGRPCServer(cfg *config.GRPCConfig, user api.UserService, secret api.SecretService) *grpc.Server {
	server := grpc.NewServer()

	api.Register(server, user, secret)

	return server
}
