package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/oscaralfaguz47/issuebot/internal/adapter/memory"
	"github.com/oscaralfaguz47/issuebot/internal/usecase"
)

func TestCreateProject_OK(t *testing.T) {
	// arrange: fake en memoria + use case
	repo := memory.NewProjectRepo()
	membershipRepo := memory.NewMembershipRepo("owner", nil) // rol owner, sin error
	uc := usecase.NewCreateProjectUseCase(repo, membershipRepo)

	err := uc.Execute(context.Background(), "user-1", "org-1", "Mi proyecto")

	// assert: no debería haber error
	if err != nil {
		t.Fatalf("esperaba nil, obtuve error: %v", err)
	}

	if repo.Count() != 1 {
		t.Fatalf("esperaba 1 proyecto guardado, hay %d", repo.Count())
	}
}

func TestCreateProject_Validacion(t *testing.T) {
	casos := []struct {
		caseName      string
		orgID         string
		name          string
		expectedError error
	}{
		{"org vacío", "", "Proyecto", usecase.ErrOrgIDRequired},
		{"name vacío", "org-1", "", usecase.ErrNameRequired},
	}

	for _, c := range casos {
		t.Run(c.caseName, func(t *testing.T) {
			repo := memory.NewProjectRepo()
			membershipRepo := memory.NewMembershipRepo("owner", nil) // rol owner, sin error
			uc := usecase.NewCreateProjectUseCase(repo, membershipRepo)

			err := uc.Execute(context.Background(), "user-1", c.orgID, c.name)

			if !errors.Is(err, c.expectedError) {
				t.Fatalf("esperaba %v, obtuve %v", c.expectedError, err)
			}
		})
	}
}

func TestCreateProject_ViewerForbidden(t *testing.T) {
	repo := memory.NewProjectRepo()
	membershipRepo := memory.NewMembershipRepo("viewer", nil) // rol viewer
	uc := usecase.NewCreateProjectUseCase(repo, membershipRepo)

	err := uc.Execute(context.Background(), "user-1", "org-1", "Mi proyecto")

	if !errors.Is(err, usecase.ErrForbidden) {
		t.Fatalf("esperaba ErrForbidden, obtuve %v", err)
	}
	if repo.Count() != 0 {
		t.Fatalf("un viewer no debería guardar nada, hay %d", repo.Count())
	}
}
