// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSecret = `-- name: CreateSecret :one
INSERT INTO secrets (name, type, metadata, data, comment, user_id)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
`

type CreateSecretParams struct {
	Name     string
	Type     SecretType
	Metadata []byte
	Data     []byte
	Comment  *string
	UserID   uuid.UUID
}

func (q *Queries) CreateSecret(ctx context.Context, arg CreateSecretParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createSecret,
		arg.Name,
		arg.Type,
		arg.Metadata,
		arg.Data,
		arg.Comment,
		arg.UserID,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (login, password)
VALUES ($1, $2) RETURNING id
`

type CreateUserParams struct {
	Login    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Login, arg.Password)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteSecret = `-- name: DeleteSecret :exec
DELETE
FROM secrets
WHERE id = $1
`

func (q *Queries) DeleteSecret(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteSecret, id)
	return err
}

const getSecret = `-- name: GetSecret :one
SELECT id, name, type, metadata, data, comment, created_at, updated_at, user_id
FROM secrets
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSecret(ctx context.Context, id uuid.UUID) (Secret, error) {
	row := q.db.QueryRow(ctx, getSecret, id)
	var i Secret
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Metadata,
		&i.Data,
		&i.Comment,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
	)
	return i, err
}

const getSecretForUpdate = `-- name: GetSecretForUpdate :one
SELECT id, name, type, metadata, data, comment, created_at, updated_at, user_id
FROM secrets
WHERE id = $1
LIMIT 1
FOR UPDATE
`

func (q *Queries) GetSecretForUpdate(ctx context.Context, id uuid.UUID) (Secret, error) {
	row := q.db.QueryRow(ctx, getSecretForUpdate, id)
	var i Secret
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Metadata,
		&i.Data,
		&i.Comment,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, login, password, created_at
FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Login,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const listSecrets = `-- name: ListSecrets :many
SELECT id, name, type, metadata, comment, created_at, updated_at
FROM secrets
WHERE user_id = $1
ORDER BY updated_at DESC
`

type ListSecretsRow struct {
	ID        uuid.UUID
	Name      string
	Type      SecretType
	Metadata  []byte
	Comment   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) ListSecrets(ctx context.Context, userID uuid.UUID) ([]ListSecretsRow, error) {
	rows, err := q.db.Query(ctx, listSecrets, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListSecretsRow
	for rows.Next() {
		var i ListSecretsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.Metadata,
			&i.Comment,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSecret = `-- name: UpdateSecret :exec
UPDATE secrets
SET name = $2, type = $3, metadata = $4, data = $5, comment = $6, created_at = $7, updated_at = $8, user_id = $9
WHERE id = $1
`

type UpdateSecretParams struct {
	ID        uuid.UUID
	Name      string
	Type      SecretType
	Metadata  []byte
	Data      []byte
	Comment   *string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
}

func (q *Queries) UpdateSecret(ctx context.Context, arg UpdateSecretParams) error {
	_, err := q.db.Exec(ctx, updateSecret,
		arg.ID,
		arg.Name,
		arg.Type,
		arg.Metadata,
		arg.Data,
		arg.Comment,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
	)
	return err
}
