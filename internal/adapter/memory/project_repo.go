package memory

import (
	"context"

	"github.com/oscaralfaguz47/issuebot/internal/domain"
)

type ProjectRepo struct {
	data map[string]domain.Project // ID -> projectID -> Project
}

func NewProjectRepo() *ProjectRepo {
	return &ProjectRepo{
		data: make(map[string]domain.Project),
	}
}

func (r *ProjectRepo) Save(ctx context.Context, project domain.Project) error {
	r.data[project.ID] = project
	return nil
}

func (r *ProjectRepo) Count() int {
	return len(r.data)
}

func (r *ProjectRepo) FindByInstallationID(ctx context.Context, installationID string) (*domain.Project, error) {
	for _, p := range r.data {
		if p.GitHubInstallationID == installationID {
			return &p, nil
		}
	}
	return nil, domain.ErrProjectNotFound
}
