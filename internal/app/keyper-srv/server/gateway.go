package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/RomanAgaltsev/keyper/internal/logger/sl"
	pb "github.com/RomanAgaltsev/keyper/pkg/keyper/v1"
)

func NewGatewayServer(log *slog.Logger, grpcAddr string, gatewayAddr string) *http.Server {
	const op = "server.NewGatewayServer"

	log = log.With(slog.String("op", op))

	conn, err := grpc.NewClient(grpcAddr)
	if err != nil {
		log.Error("new gRPC client", sl.Err(err))
	}

	mux := runtime.NewServeMux()
	if err := pb.RegisterUserServiceHandler(context.Background(), mux, conn); err != nil {
		log.Error("register user service handler", sl.Err(err))
	}

	if err := pb.RegisterSecretServiceHandler(context.Background(), mux, conn); err != nil {
		log.Error("register secret service handler", sl.Err(err))
	}

	gatewayServer := &http.Server{
		Addr:    gatewayAddr,
		Handler: mux,
	}

	return gatewayServer
}
