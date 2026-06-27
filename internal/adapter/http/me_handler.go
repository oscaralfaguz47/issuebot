package http

import (
	"encoding/json"
	"net/http"

	"github.com/oscaralfaguz47/issuebot/internal/domain"
)

func HandleMe(membershipRepo domain.MembershipRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "no user in context", http.StatusUnauthorized)
			return
		}

		memberships, err := membershipRepo.ListByUser(r.Context(), userID)
		if err != nil {
			http.Error(w, "could not load memberships", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(memberships)
	}
}
