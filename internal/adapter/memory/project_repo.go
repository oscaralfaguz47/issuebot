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
