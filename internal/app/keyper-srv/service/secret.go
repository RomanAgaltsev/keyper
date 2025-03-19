package service

import (
	"context"
	"log/slog"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/api"
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

func (s *SecretService) Create(ctx context.Context) error {
	return nil
}

func (s *SecretService) Get(ctx context.Context) error {
	return nil
}

func (s *SecretService) List(ctx context.Context) error {
	return nil
}

func (s *SecretService) Update(ctx context.Context) error {
	return nil
}

func (s *SecretService) Delete(ctx context.Context) error {
	return nil
}
