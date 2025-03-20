package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WithTxFunc func(ctx context.Context, tx pgx.Tx) error

func WithTx(ctx context.Context, db *pgxpool.Pool, fn WithTxFunc) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db.Begin(ctx): %w", err)
	}

	if err = fn(ctx, tx); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			return fmt.Errorf("Tx.Rollback(ctx): %w", err)
		}

		return fmt.Errorf("Tx.WithTxFunc: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("Tx.Commit: %w", err)
	}

	return nil
}
