package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/cenkalti/backoff/v5"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/repository"
	"github.com/RomanAgaltsev/keyper/internal/model"
)

var (
	_ UserRepository = (*repository.UserRepository)(nil)

	ErrLoginTaken = errors.New("login has already been taken")
	ErrLoginWrong = errors.New("wrong login/password")
)

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
	// TODO: password hashing

	// Create user in the repository
	err := s.repository.Create(ctx, repository.DefaultRetryOpts, user)

	// There is a conflict - the login is already exists in the database
	if errors.Is(err, repository.ErrConflict) {
		return ErrLoginTaken
	}

	// There is another error
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, user model.User) error {
	// Ger user from repository
	userInRepo, err := s.repository.Get(ctx, repository.DefaultRetryOpts, user.Login)
	if err != nil {
		return err
	}

	// TODO: password hashing

	// If user doesn`t exist or password is wrong
	if (userInRepo == model.User{}) || !(user.Password == userInRepo.Password) {
		return ErrLoginWrong
	}

	return nil
}
