package service

import (
	"context"
	"log/slog"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-cli/client"
	"github.com/RomanAgaltsev/keyper/internal/model"
)

var _ UserClient = (*client.UserClient)(nil)

type UserClient interface {
	Register(ctx context.Context, user *model.User) error
	Login(ctx context.Context, user *model.User) error
}

func NewUserService(log *slog.Logger, client UserClient) *UserService {
	return &UserService{
		log:    log,
		client: client,
	}
}

type UserService struct {
	log *slog.Logger

	client UserClient
}

func (s *UserService) Register(ctx context.Context, user *model.User) error {
	return nil
}

func (s *UserService) Login(ctx context.Context, user *model.User) error {
	return nil
}
