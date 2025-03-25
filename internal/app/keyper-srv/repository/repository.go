package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/RomanAgaltsev/keyper/internal/database"
	"github.com/RomanAgaltsev/keyper/internal/database/queries"
	"github.com/RomanAgaltsev/keyper/internal/model"
	"github.com/RomanAgaltsev/keyper/pkg/transform"
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

func (r *UserRepository) Create(
	ctx context.Context,
	ro []backoff.RetryOption,
	user *model.User,
) error {
	// PG error to catch the conflict
	var pgErr *pgconn.PgError

	createUserParams := transform.UserToCreateUserParams(user)

	// Create a function to wrap user creation with exponential backoff
	f := func() (error, error) {
		// Create user
		_, err := r.q.CreateUser(ctx, createUserParams)

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

func (r *UserRepository) Get(
	ctx context.Context,
	ro []backoff.RetryOption,
	userID uuid.UUID,
) (
	*model.User,
	error,
) {
	// Get user from DB
	userDB, err := backoff.Retry(ctx, func() (queries.User, error) {
		return r.q.GetUser(ctx, userID)
	}, ro...)

	// Check if something has gone wrong
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Check if there is nothing to return
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	user := transform.DBToUser(userDB)

	// Return user
	return user, nil
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

func (r *SecretRepository) Create(
	ctx context.Context,
	ro []backoff.RetryOption,
	secret *model.Secret,
) (
	uuid.UUID,
	error,
) {
	createSecretParams := transform.SecretToCreateSecretParams(secret)

	secretID, err := backoff.Retry(ctx, func() (uuid.UUID, error) {
		return r.q.CreateSecret(ctx, createSecretParams)
	}, ro...)

	if err != nil {
		return uuid.Nil, err
	}

	return secretID, nil
}

func (r *SecretRepository) Get(
	ctx context.Context,
	ro []backoff.RetryOption,
	secretID uuid.UUID,
) (
	*model.Secret,
	error,
) {
	secretDB, err := backoff.Retry(ctx, func() (queries.Secret, error) {
		return r.q.GetSecret(ctx, secretID)
	}, ro...)

	if err != nil {
		return nil, err
	}

	secret := transform.DBToSecret(secretDB)

	return secret, nil
}

func (r *SecretRepository) List(
	ctx context.Context,
	ro []backoff.RetryOption,
	userID uuid.UUID,
) (
	model.Secrets,
	error,
) {
	listSecretsRow, err := backoff.Retry(ctx, func() ([]queries.ListSecretsRow, error) {
		return r.q.ListSecrets(ctx, userID)
	}, ro...)

	if err != nil {
		return nil, err
	}

	secrets := transform.ListSecretsRowToSecrets(listSecretsRow)

	return secrets, nil
}

func (r *SecretRepository) Update(
	ctx context.Context,
	ro []backoff.RetryOption,
	secret *model.Secret,
	updateFn func(secretTo, secretFrom *model.Secret) (bool, error),
) error {
	return database.WithTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		// Create query with transaction
		qtx := r.q.WithTx(tx)

		secretDB, err := backoff.Retry(ctx, func() (queries.Secret, error) {
			return qtx.GetSecretForUpdate(ctx, secret.ID)
		}, ro...)

		if err != nil {
			return err
		}

		secretTo := transform.DBToSecret(secretDB)

		ok, err := updateFn(secretTo, secret)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		updateSecretParams := transform.SecretToUpdateSecretParams(secretTo)

		_, err = backoff.Retry(ctx, func() (bool, error) {
			err := qtx.UpdateSecret(ctx, updateSecretParams)
			if err != nil {
				return false, err
			}
			return true, nil
		}, ro...)

		if err != nil {
			return err
		}

		return nil
	})
}

func (r *SecretRepository) Delete(
	ctx context.Context,
	ro []backoff.RetryOption,
	secretID uuid.UUID,
) error {
	_, err := backoff.Retry(ctx, func() (bool, error) {
		err := r.q.DeleteSecret(ctx, secretID)
		if err != nil {
			return false, err
		}
		return true, nil
	}, ro...)

	if err != nil {
		return err
	}

	return nil
}
