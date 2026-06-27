package domain

import (
	"context"
	"errors"
)

var ErrNotMember = errors.New("user is not a member of this org")

type MembershipRepository interface {
	GetRole(ctx context.Context, orgID, userID string) (string, error)
	ListByUser(ctx context.Context, userID string) ([]Membership, error)
}
type Membership struct {
	OrgID  string `json:"org_id"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}
