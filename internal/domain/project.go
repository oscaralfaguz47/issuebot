package domain

import (
	"context"
	"time"
)

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
}
