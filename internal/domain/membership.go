package domain

import (
	"context"
	"errors"
)

var ErrNotMember = errors.New("user is not a member of this org")

type MembershipRepository interface {
	GetRole(ctx context.Context, orgID, userID string) (string, error)
}