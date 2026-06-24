package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/oscaralfaguz47/issuebot/internal/domain"
	"github.com/oscaralfaguz47/issuebot/internal/usecase"
)

func HandleCreateProject(uc *usecase.CreateProjectUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			OrgID string `json:"org_id"`
			Name  string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		userID, ok := UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "no user in context", http.StatusUnauthorized)
			return
		}

		if err := uc.Execute(r.Context(), userID, input.OrgID, input.Name); err != nil {
			if errors.Is(err, usecase.ErrForbidden) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "Project created successfully")
	}
}

func HandleListProjects(repo domain.ProjectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "no user in context", http.StatusUnauthorized)
			return
		}

		projects, err := repo.ListByUser(r.Context(), userID)
		if err != nil {
			http.Error(w, "could not list projects", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(projects)
	}
}
