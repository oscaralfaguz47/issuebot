package postgres

import (
	"context"

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
