package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/RomanAgaltsev/keyper/internal/database/queries"
	"github.com/RomanAgaltsev/keyper/internal/model"
)

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

func (r *UserRepository) Create(ctx context.Context, user model.User) error {
	return nil
}

func (r *UserRepository) Get(ctx context.Context, login string) (model.User, error) {
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

func (r *SecretRepository) Create(ctx context.Context, secret model.Secret) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (r *SecretRepository) Get(ctx context.Context, secretID uuid.UUID) (model.Secret, error) {
	return model.Secret{}, nil
}

func (r *SecretRepository) List(ctx context.Context, user model.User) (model.Secrets, error) {
	return nil, nil
}

func (r *SecretRepository) Update(ctx context.Context, secret model.Secret) error {
	return nil
}

func (r *SecretRepository) Delete(ctx context.Context, secretID uuid.UUID) error {
	return nil
}
