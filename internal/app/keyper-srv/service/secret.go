package service

import (
	"context"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/api"
	"github.com/RomanAgaltsev/keyper/internal/config"
)

var _ api.SecretService = (*SecretService)(nil)

func NewSecretService(cfg *config.AppConfig) *SecretService {
	return &SecretService{}
}

type SecretService struct{}

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
