package usecase

import (
	"context"
	"errors"

	"github.com/oscaralfaguz47/issuebot/internal/domain"
)

var (
	ErrOrgIDRequired = errors.New("org_id is required")
	ErrNameRequired  = errors.New("name is required")
	ErrForbidden     = errors.New("user not allowed to perform this action")
)

type CreateProjectUseCase struct {
	repo           domain.ProjectRepository
	membershipRepo domain.MembershipRepository
}

func NewCreateProjectUseCase(repo domain.ProjectRepository, membershipRepo domain.MembershipRepository) *CreateProjectUseCase {
	return &CreateProjectUseCase{repo: repo, membershipRepo: membershipRepo}
}

func (uc *CreateProjectUseCase) Execute(ctx context.Context, userID string, orgID string, name string) error {
	if orgID == "" {
		return ErrOrgIDRequired
	}
	if name == "" {
		return ErrNameRequired
	}
	role, err := uc.membershipRepo.GetRole(ctx, orgID, userID)
	if err != nil {
		return ErrForbidden
	}
	if role == "viewer" {
		return ErrForbidden
	}
	return uc.repo.Save(ctx, domain.Project{
		OrgID: orgID,
		Name:  name,
	})
}
