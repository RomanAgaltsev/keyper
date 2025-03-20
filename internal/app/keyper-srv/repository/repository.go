package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/RomanAgaltsev/keyper/internal/database/queries"
	"github.com/RomanAgaltsev/keyper/internal/model"
)

const DefaultRetryMaxElapsedTime = 5 * time.Second

var (
	ErrConflict = errors.New("data conflict")

	DefaultRetryOpts = []backoff.RetryOption{
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxElapsedTime(DefaultRetryMaxElapsedTime),
	}
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

func (r *UserRepository) Create(ctx context.Context, ro []backoff.RetryOption, user model.User) error {
	// PG error to catch the conflict
	var pgErr *pgconn.PgError

	// Create a function to wrap user creation with exponential backoff
	f := func() (error, error) {
		// Create user
		_, err := r.q.CreateUser(ctx, queries.CreateUserParams{
			Login:    user.Login,
			Password: user.Password,
		})

		// Check if there is a conflict
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return ErrConflict, nil
		}

		// Check if something has gone wrong
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Call the wrapping function
	errConf, err := backoff.Retry(ctx, f, ro...)
	if err != nil {
		return err
	}

	// There is a conflict
	if errConf != nil {
		return errConf
	}

	return nil
}

func (r *UserRepository) Get(ctx context.Context, ro []backoff.RetryOption, login string) (model.User, error) {
	// Create a function to wrap user getting with exponential backoff
	f := func() (queries.GetUserRow, error) {
		return r.q.GetUser(ctx, login)
	}

	var user model.User

	// Get user from DB
	userRow, err := backoff.Retry(ctx, f, ro...)

	// Check if something has gone wrong
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return user, err
	}

	// Check if there is nothing to return
	if errors.Is(err, sql.ErrNoRows) {
		return user, nil
	}

	// Return user
	return model.User{
		Login:    login,
		Password: userRow.Password,
	}, nil
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
