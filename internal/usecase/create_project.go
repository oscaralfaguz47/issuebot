package usecase

import (
	"context"

	"github.com/oscaralfaguz47/issuebot/internal/domain"
)

type CreateProjectUseCase struct {
	repo domain.ProjectRepository
}

func NewCreateProjectUseCase(repo domain.ProjectRepository) *CreateProjectUseCase {
	return &CreateProjectUseCase{repo: repo}
}

func (uc *CreateProjectUseCase) Execute(ctx context.Context, orgID string, name string) error {
	return uc.repo.Save(ctx, domain.Project{
		OrgID: orgID,
		Name:  name,
	})
}
