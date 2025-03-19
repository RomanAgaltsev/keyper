package service

import (
	"context"
	"log/slog"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/api"
)

var _ api.UserService = (*UserService)(nil)

func NewUserService(log *slog.Logger) *UserService {
	return &UserService{
		log: log,
	}
}

type UserService struct {
	log *slog.Logger
}

func (s *UserService) Register(ctx context.Context) error {
	return nil
}

func (s *UserService) Login(ctx context.Context) error {
	return nil
}
