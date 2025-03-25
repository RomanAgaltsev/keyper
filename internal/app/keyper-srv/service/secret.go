package service

import (
	"bytes"
	"context"
	"log/slog"

	"github.com/cenkalti/backoff/v5"
	"github.com/google/uuid"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/repository"
	"github.com/RomanAgaltsev/keyper/internal/model"
)

const dataPortionSize = 1024 * 1024

var _ SecretRepository = (*repository.SecretRepository)(nil)

type SecretRepository interface {
	Create(ctx context.Context, ro []backoff.RetryOption, secret *model.Secret) (uuid.UUID, error)
	Update(ctx context.Context, ro []backoff.RetryOption, userID uuid.UUID, secret *model.Secret, updateFn func(dst, src *model.Secret) (bool, error)) error
	Get(ctx context.Context, ro []backoff.RetryOption, userID uuid.UUID, secretID uuid.UUID) (*model.Secret, error)
	List(ctx context.Context, ro []backoff.RetryOption, userID uuid.UUID) (model.Secrets, error)
	Delete(ctx context.Context, ro []backoff.RetryOption, userID uuid.UUID, secretID uuid.UUID) error
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

func (s *SecretService) Create(ctx context.Context, secret *model.Secret) (uuid.UUID, error) {
	secretID, err := s.repository.Create(ctx, repository.DefaultRetryOpts, secret)
	if err != nil {
		return uuid.Nil, err
	}

	return secretID, nil
}

func (s *SecretService) Update(ctx context.Context, userID uuid.UUID, secret *model.Secret) error {
	return s.repository.Update(ctx, repository.DefaultRetryOpts, userID, secret, func(secretTo, secretFrom *model.Secret) (bool, error) {
		err := secretTo.UpdateWith(secretFrom)
		if err != nil {
			return false, nil
		}
		return true, nil
	})
}

func (s *SecretService) UpdateData(ctx context.Context, userID uuid.UUID, secretID uuid.UUID, dataCh chan []byte) error {
	secret, err := s.repository.Get(ctx, repository.DefaultRetryOpts, userID, secretID)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	for dataChunk := range dataCh {
		if buf.Len()+len(dataChunk) > dataPortionSize {
			// TODO: save portion of data and clear dataPortion

			buf.Reset()
		}
		buf.Write(dataChunk)
	}

	if buf.Len() > 0 {
		// TODO: handle the remainder
	}

	return nil
}

func (s *SecretService) Get(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) (*model.Secret, error) {
	return s.repository.Get(ctx, repository.DefaultRetryOpts, userID, secretID)
}

func (s *SecretService) GetData(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) (chan []byte, error) {
	return nil, nil
}

func (s *SecretService) List(ctx context.Context, userID uuid.UUID) (model.Secrets, error) {
	return s.repository.List(ctx, repository.DefaultRetryOpts, userID)
}

func (s *SecretService) Delete(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) error {
	return s.repository.Delete(ctx, repository.DefaultRetryOpts, userID, secretID)
}
