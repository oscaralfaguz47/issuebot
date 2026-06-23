package memory

import (
	"context"
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
