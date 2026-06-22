package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oscaralfaguz47/issuebot/internal/domain"
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

func (r *JobRepo) ClaimJob(ctx context.Context) (*domain.Job, error) {
	query := `
		update jobs
		set status = 'processing', updated_at = now()
		where id = (
			select id from jobs
			where status = 'pending'
			order by created_at
			for update skip locked
			limit 1
		)
		returning id, project_id, type, payload, attempts
	`
	row := r.pool.QueryRow(ctx, query)

	var job domain.Job
	err := row.Scan(&job.ID, &job.ProjectID, &job.Type, &job.Payload, &job.Attempts)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNoJobs
		}
		return nil, err
	}
	return &job, nil
}

func (r *JobRepo) MarkDone(ctx context.Context, jobID string) error {
	_, err := r.pool.Exec(ctx, `update jobs set status = 'done', updated_at = now() where id = $1`, jobID)
	return err
}

func (r *JobRepo) MarkFailed(ctx context.Context, jobID string) error {
	_, err := r.pool.Exec(ctx, `update jobs set status = 'failed', updated_at = now() where id = $1`, jobID)
	return err
}

func (r *JobRepo) Reschedule(ctx context.Context, jobID string) error {
	_, err := r.pool.Exec(ctx, `update jobs set status = 'pending', attempts = attempts + 1, updated_at = now() where id = $1`, jobID)
	return err
}
