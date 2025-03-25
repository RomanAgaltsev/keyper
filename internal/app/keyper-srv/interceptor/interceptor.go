package interceptor

import (
	"context"
	"log/slog"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	"github.com/RomanAgaltsev/keyper/internal/pkg/auth"
	pb "github.com/RomanAgaltsev/keyper/pkg/keyper/v1"
)

var (
	recoveryHandlerFunc recovery.RecoveryHandlerFunc

	methodsRequireAuth map[string]bool = map[string]bool{
		pb.UserService_RegisterUserV1_FullMethodName:   false,
		pb.UserService_LoginUserV1_FullMethodName:      false,
		pb.SecretService_CreateSecretV1_FullMethodName: true,
		pb.SecretService_GetSecretV1_FullMethodName:    true,
		pb.SecretService_ListSecretsV1_FullMethodName:  true,
		pb.SecretService_UpdateSecretV1_FullMethodName: true,
		pb.SecretService_DeleteSecretV1_FullMethodName: true,
	}
)

func Logger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func LoggerOpts() []logging.Option {
	return []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}
}

func RecoveryOpts() []recovery.Option {
	recoveryHandlerFunc = func(p any) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}

	return []recovery.Option{
		recovery.WithRecoveryHandler(recoveryHandlerFunc),
	}
}

func NewAuthInterceptor(secretKey string) func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if _, ok := methodsRequireAuth[info.FullMethod]; !ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "missing metadata")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
		if tokenString == "" {
			return nil, status.Error(codes.Unauthenticated, "malformed token")
		}

		ja := auth.NewAuth(secretKey)

		token, err := ja.Decode(tokenString)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		if err = jwt.Validate(token, ja.ValidateOptions()...); err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		// TODO: check toket lifetime

		claims := token.PrivateClaims()

		ctx = context.WithValue(ctx, auth.UserIDClaimName, claims[string(auth.UserIDClaimName)])

		return handler(ctx, req)
	}
}
