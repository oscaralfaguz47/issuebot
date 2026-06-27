package memory

import (
	"context"

	"github.com/oscaralfaguz47/issuebot/internal/domain"
)

type MembershipRepo struct {
	role string
	err  error
}

func NewMembershipRepo(role string, err error) *MembershipRepo {
	return &MembershipRepo{role: role, err: err}
}

func (r *MembershipRepo) GetRole(ctx context.Context, orgID, userID string) (string, error) {
	return r.role, r.err
}

func (r *MembershipRepo) ListByUser(ctx context.Context, userID string) ([]domain.Membership, error) {
	if r.err != nil {
		return nil, r.err
	}
	return []domain.Membership{
		{OrgID: "org-1", UserID: userID, Role: r.role},
	}, nil
}