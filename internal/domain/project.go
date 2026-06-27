package domain

import (
	"context"
	"errors"
	"time"
)

var ErrProjectNotFound = errors.New("project not found")

type Project struct {
	ID                   string    `json:"id"`
	OrgID                string    `json:"org_id"`
	Name                 string    `json:"name"`
	GitHubRepo           string    `json:"github_repo"`
	GitHubInstallationID string    `json:"github_installation_id"`
	CreatedAt            time.Time `json:"created_at"`
}

type ProjectRepository interface {
	Save(ctx context.Context, project Project) error
	FindByInstallationID(ctx context.Context, installationID string) (*Project, error)
	ListByUser(ctx context.Context, userID string) ([]Project, error)
}
