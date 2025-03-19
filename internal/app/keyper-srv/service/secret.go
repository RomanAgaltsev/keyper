package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/api"
	"github.com/RomanAgaltsev/keyper/internal/model"
)

var _ api.SecretService = (*SecretService)(nil)

func NewSecretService(log *slog.Logger) *SecretService {
	return &SecretService{
		log: log,
	}
}

type SecretService struct {
	log *slog.Logger
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
