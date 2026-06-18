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
	uc := usecase.NewCreateProjectUseCase(repo)

	// act: ejecuto con input válido
	err := uc.Execute(context.Background(), "org-1", "Mi proyecto")

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
			uc := usecase.NewCreateProjectUseCase(repo)

			err := uc.Execute(context.Background(), c.orgID, c.name)

			if !errors.Is(err, c.expectedError) {
				t.Fatalf("esperaba %v, obtuve %v", c.expectedError, err)
			}
		})
	}
}
