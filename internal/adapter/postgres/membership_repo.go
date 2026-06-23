package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oscaralfaguz47/issuebot/internal/domain"
)

type MembershipRepo struct {
	pool *pgxpool.Pool
}

func NewMembershipRepo(pool *pgxpool.Pool) *MembershipRepo {
	return &MembershipRepo{pool: pool}
}

func (r *MembershipRepo) GetRole(ctx context.Context, orgID, userID string) (string, error) {
	query := `
		select role from memberships
		where org_id = $1 and user_id = $2
	`
	row := r.pool.QueryRow(ctx, query, orgID, userID)

	var role string
	err := row.Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrNotMember
		}
		return "", err
	}
	return role, nil
}

