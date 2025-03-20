package service

import (
	"context"
	"log/slog"

	"github.com/cenkalti/backoff/v5"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/repository"
	"github.com/RomanAgaltsev/keyper/internal/model"
)

var _ UserRepository = (*repository.UserRepository)(nil)

type UserRepository interface {
	Create(ctx context.Context, ro []backoff.RetryOption, user model.User) error
	Get(ctx context.Context, ro []backoff.RetryOption, login string) (model.User, error)
}

func NewUserService(log *slog.Logger, repository *repository.UserRepository) *UserService {
	return &UserService{
		log:        log,
		repository: repository,
	}
}

type UserService struct {
	log *slog.Logger

	repository *repository.UserRepository
}

func (s *UserService) Register(ctx context.Context, user model.User) error {
	return nil
}

func (s *UserService) Login(ctx context.Context, user model.User) error {
	return nil
}
