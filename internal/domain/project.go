package domain

import (
	"context"
	"errors"
	"time"
)

var ErrProjectNotFound = errors.New("project not found")

type Project struct {
	ID                   string
	OrgID                string
	Name                 string
	GitHubRepo           string
	GitHubInstallationID string
	CreatedAt            time.Time
}

type ProjectRepository interface {
	Save(ctx context.Context, project Project) error
	FindByInstallationID(ctx context.Context, installationID string) (*Project, error)
	ListByUser(ctx context.Context, userID string) ([]Project, error)
}
