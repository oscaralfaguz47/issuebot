package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/oscaralfaguz47/issuebot/internal/domain"
)

type ProjectRepo struct {
	pool *pgxpool.Pool
}

func NewProjectRepo(pool *pgxpool.Pool) *ProjectRepo {
	return &ProjectRepo{pool: pool}
}

func (r *ProjectRepo) Save(ctx context.Context, project domain.Project) error {
	query := `
		insert into projects (org_id, name, github_repo, github_installation_id)
		values ($1, $2, $3, $4)
	`
	_, err := r.pool.Exec(ctx, query, project.OrgID, project.Name, project.GitHubRepo, project.GitHubInstallationID)
	return err
}

func (r *ProjectRepo) FindByInstallationID(ctx context.Context, installationID string) (*domain.Project, error) {
	query := `
		select id, org_id, name, github_repo, github_installation_id, created_at
		from projects
		where github_installation_id = $1
	`
	row := r.pool.QueryRow(ctx, query, installationID)

	var p domain.Project
	err := row.Scan(&p.ID, &p.OrgID, &p.Name, &p.GitHubRepo, &p.GitHubInstallationID, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProjectNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepo) ListByUser(ctx context.Context, userID string) ([]domain.Project, error) {
	query := `select p.id, p.org_id, p.name, p.github_repo, p.github_installation_id, p.created_at
from projects p
join memberships m on m.org_id = p.org_id
where m.user_id = $1`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(&p.ID, &p.OrgID, &p.Name, &p.GitHubRepo, &p.GitHubInstallationID, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}
