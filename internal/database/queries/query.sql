-- name: CreateUser :one
INSERT INTO users (login, password)
VALUES ($1, $2) RETURNING id;

-- name: GetUser :one
SELECT id, login, password, created_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateSecret :one
INSERT INTO secrets (name, type, metadata, data, comment, user_id)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;

-- name: GetSecret :one
SELECT id, name, type, metadata, data, comment, created_at, updated_at, user_id
FROM secrets
WHERE id = $1 LIMIT 1;

-- name: GetSecretForUpdate :one
SELECT id, name, type, metadata, data, comment, created_at, updated_at, user_id
FROM secrets
WHERE id = $1
LIMIT 1
FOR UPDATE;

-- name: ListSecrets :many
SELECT id, name, type, metadata, comment, created_at, updated_at
FROM secrets
WHERE user_id = $1
ORDER BY updated_at DESC;

-- name: UpdateSecret :exec
UPDATE secrets
SET name = $2, type = $3, metadata = $4, data = $5, comment = $6, created_at = $7, updated_at = $8, user_id = $9
WHERE id = $1;

-- name: DeleteSecret :exec
DELETE
FROM secrets
WHERE id = $1;