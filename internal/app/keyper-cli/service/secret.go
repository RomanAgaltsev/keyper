package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-cli/client"
	"github.com/RomanAgaltsev/keyper/internal/model"
)

var _ SecretClient = (*client.SecretClient)(nil)

type SecretClient interface {
	Create(ctx context.Context, secret *model.Secret) (uuid.UUID, error)
	Update(ctx context.Context, userID uuid.UUID, secret *model.Secret) error
	UpdateData(ctx context.Context, userID uuid.UUID, secret *model.Secret) error
	Get(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) (*model.Secret, error)
	GetData(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) (*model.Secret, error)
	List(ctx context.Context, userID uuid.UUID) (model.Secrets, error)
	Delete(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) error
}

func NewSecretService(log *slog.Logger, client SecretClient) *SecretService {
	return &SecretService{
		log:    log,
		client: client,
	}
}

type SecretService struct {
	log *slog.Logger

	client SecretClient
}
