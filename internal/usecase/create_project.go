package usecase

import (
	"context"
	"errors"

	"github.com/oscaralfaguz47/issuebot/internal/domain"
)

var (
	ErrOrgIDRequired = errors.New("org_id is required")
	ErrNameRequired  = errors.New("name is required")
)

type CreateProjectUseCase struct {
	repo domain.ProjectRepository
}

func NewCreateProjectUseCase(repo domain.ProjectRepository) *CreateProjectUseCase {
	return &CreateProjectUseCase{repo: repo}
}

func (uc *CreateProjectUseCase) Execute(ctx context.Context, orgID string, name string) error {
	if orgID == "" {
		return ErrOrgIDRequired
	}
	if name == "" {
		return ErrNameRequired
	}
	return uc.repo.Save(ctx, domain.Project{
		OrgID: orgID,
		Name:  name,
	})
}
