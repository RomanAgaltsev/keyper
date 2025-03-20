package repository

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/RomanAgaltsev/keyper/internal/database/queries"
	"github.com/RomanAgaltsev/keyper/internal/model"
)

const DefaultRetryMaxElapsedTime = 5 * time.Second

var DefaultRetryOpts = []backoff.RetryOption{
	backoff.WithBackOff(backoff.NewExponentialBackOff()),
	backoff.WithMaxElapsedTime(DefaultRetryMaxElapsedTime),
}

func NewUserRepository(dbpool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: dbpool,
		q:  queries.New(dbpool),
	}
}

type UserRepository struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

func (r *UserRepository) Create(ctx context.Context, ro []backoff.RetryOption, user model.User) error {
	return nil
}

func (r *UserRepository) Get(ctx context.Context, ro []backoff.RetryOption, login string) (model.User, error) {
	return model.User{}, nil
}

func NewSecretRepository(dbpool *pgxpool.Pool) *SecretRepository {
	return &SecretRepository{
		db: dbpool,
		q:  queries.New(dbpool),
	}
}

type SecretRepository struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

func (r *SecretRepository) Create(ctx context.Context, ro []backoff.RetryOption, secret model.Secret) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (r *SecretRepository) Get(ctx context.Context, ro []backoff.RetryOption, secretID uuid.UUID) (model.Secret, error) {
	return model.Secret{}, nil
}

func (r *SecretRepository) List(ctx context.Context, ro []backoff.RetryOption, user model.User) (model.Secrets, error) {
	return nil, nil
}

func (r *SecretRepository) Update(ctx context.Context, ro []backoff.RetryOption, secret model.Secret) error {
	return nil
}

func (r *SecretRepository) Delete(ctx context.Context, ro []backoff.RetryOption, secretID uuid.UUID) error {
	return nil
}
