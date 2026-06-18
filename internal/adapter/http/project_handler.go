package http

import (
	"encoding/json"
	"fmt"
	"net/http"

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

		if err := uc.Execute(r.Context(), input.OrgID, input.Name); err != nil {
			http.Error(w, "...", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "Project created successfully")
	}
}
