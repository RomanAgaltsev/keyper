-- +goose Up
CREATE TABLE users (
    id         UUID               NOT NULL DEFAULT uuid_generate_v4(),
    login      VARCHAR(30) UNIQUE NOT NULL,
    password   VARCHAR(60)        NOT NULL,
    created_at TIMESTAMP          NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;