-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE secret_type AS ENUM ('UNSPECIFIED', 'CREDENTIALS', 'TEXT', 'BINARY', 'CARD');

CREATE TABLE secrets (
    id         UUID         NOT NULL DEFAULT uuid_generate_v4(),
    name       VARCHAR(100) UNIQUE NOT NULL,
    type       secret_type  NOT NULL DEFAULT 'UNSPECIFIED',
    metadata   BYTEA,
    data       BYTEA,
    comment    VARCHAR(100),
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    user_id    UUID REFERENCES users (id) NOT NULL,
    UNIQUE (name, user_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE secrets;
DROP TYPE secret_type;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd