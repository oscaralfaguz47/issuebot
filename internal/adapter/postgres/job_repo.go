package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepo struct {
	pool *pgxpool.Pool
}

func NewJobRepo(pool *pgxpool.Pool) *JobRepo {
	return &JobRepo{pool: pool}
}

func (r *JobRepo) Enqueue(ctx context.Context, projectID, jobType, payload, idempotencyKey string) error {
	query := `
		insert into jobs (project_id, type, payload, idempotency_key)
		values ($1, $2, $3, $4)
	`
	_, err := r.pool.Exec(ctx, query, projectID, jobType, payload, idempotencyKey)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil
		}
		return err
	}
	return nil
}
