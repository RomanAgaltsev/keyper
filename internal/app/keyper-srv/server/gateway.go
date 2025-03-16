package server

import (
	"context"
	pb "github.com/RomanAgaltsev/keyper/pkg/keyper/v1"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func NewGatewayServer(grpcAddr string, gatewayAddr string) *http.Server {
	conn, err := grpc.NewClient(grpcAddr)
	if err != nil {

	}

	mux := runtime.NewServeMux()
	if err := pb.RegisterUserServiceHandler(context.Background(), mux, conn); err != nil {

	}

	if err := pb.RegisterUserServiceHandler(context.Background(), mux, conn); err != nil {

	}

	gatewayServer := &http.Server{
		Addr:    gatewayAddr,
		Handler: mux,
	}

	return gatewayServer
}
