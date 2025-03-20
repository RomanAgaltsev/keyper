package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewConnectionPool creates new pgx connection pool and runs migrations.
func NewConnectionPool(ctx context.Context, databaseURI string) (*pgxpool.Pool, error) {
	const op = "database.NewConnectionPool"

	// Create new connection pool
	dbpool, err := pgxpool.New(ctx, databaseURI)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Ping DB
	if err = dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dbpool, nil
}
