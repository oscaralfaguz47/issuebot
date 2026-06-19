package usecase

import (
	"context"

	"github.com/oscaralfaguz47/issuebot/internal/domain"
)

type EnqueueJobUseCase struct {
	jobRepo domain.JobRepository
}

func NewEnqueueJobUseCase(jobRepo domain.JobRepository) *EnqueueJobUseCase {
	return &EnqueueJobUseCase{jobRepo: jobRepo}
}

func (uc *EnqueueJobUseCase) Execute(ctx context.Context, projectID, jobType, payload, idempotencyKey string) error {
	return uc.jobRepo.Enqueue(ctx, projectID, jobType, payload, idempotencyKey)
}
