package platform

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewDBPool abre un pool de conexiones a Postgres usando el DSN dado.
func NewDBPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
