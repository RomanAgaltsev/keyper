package service

import (
	"context"
	"log/slog"

	"github.com/cenkalti/backoff/v5"
	"github.com/google/uuid"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/repository"
	"github.com/RomanAgaltsev/keyper/internal/model"
)

var _ SecretRepository = (*repository.SecretRepository)(nil)

type SecretRepository interface {
	Create(ctx context.Context, ro []backoff.RetryOption, secret model.Secret) (uuid.UUID, error)
	Get(ctx context.Context, ro []backoff.RetryOption, secretID uuid.UUID) (model.Secret, error)
	List(ctx context.Context, ro []backoff.RetryOption, user model.User) (model.Secrets, error)
	Update(ctx context.Context, ro []backoff.RetryOption, secret model.Secret, updateFn func(dst, src model.Secret) (bool, error)) error
	Delete(ctx context.Context, ro []backoff.RetryOption, secretID uuid.UUID) error
}

func NewSecretService(log *slog.Logger, repository *repository.SecretRepository) *SecretService {
	return &SecretService{
		log:        log,
		repository: repository,
	}
}

type SecretService struct {
	log *slog.Logger

	repository *repository.SecretRepository
}

func (s *SecretService) Create(ctx context.Context, secret model.Secret) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (s *SecretService) Get(ctx context.Context, secretID uuid.UUID) (model.Secret, error) {
	return model.Secret{}, nil
}

func (s *SecretService) List(ctx context.Context, user model.User) (model.Secrets, error) {
	return nil, nil
}

func (s *SecretService) Update(ctx context.Context, secret model.Secret) error {
	return nil
}

func (s *SecretService) Delete(ctx context.Context, secretID uuid.UUID) error {
	return nil
}
