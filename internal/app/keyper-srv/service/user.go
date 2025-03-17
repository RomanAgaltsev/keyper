package service

import (
	"context"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/api"
	"github.com/RomanAgaltsev/keyper/internal/config"
)

var _ api.UserService = (*UserService)(nil)

func NewUserService(cfg *config.AppConfig) *UserService {
	return &UserService{}
}

type UserService struct{}

func (s *UserService) Register(ctx context.Context) error {
	return nil
}

func (s *UserService) Login(ctx context.Context) error {
	return nil
}
