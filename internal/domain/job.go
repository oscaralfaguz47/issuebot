package domain

import (
	"context"
	"time"
)

type Job struct {
	ID             string
	ProjectID      string
	Type           string
	Status         string
	Payload        string
	Attempts       int
	IdempotencyKey string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type JobRepository interface {
	Enqueue(ctx context.Context, projectID, jobType, payload, idempotencyKey string) error
}
