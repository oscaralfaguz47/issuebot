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