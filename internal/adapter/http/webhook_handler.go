package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/oscaralfaguz47/issuebot/internal/domain"
	"github.com/oscaralfaguz47/issuebot/internal/usecase"
)

func HandleGitHubWebhook(webhookSecret string, projectRepo domain.ProjectRepository, enqueueJob *usecase.EnqueueJobUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Can't read request body", http.StatusBadRequest)
			return
		}
		signature := r.Header.Get("X-Hub-Signature-256")
		if signature == "" {
			http.Error(w, "Missing signature", http.StatusBadRequest)
			return
		}
		if !validSignature(body, signature, webhookSecret) {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
		deliveryID := r.Header.Get("X-GitHub-Delivery")
		eventType := r.Header.Get("X-GitHub-Event")

		var payload struct {
			Installation struct {
				ID int64 `json:"id"`
			} `json:"installation"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		installationID := strconv.FormatInt(payload.Installation.ID, 10)

		project, err := projectRepo.FindByInstallationID(r.Context(), installationID)
		if err != nil {
			if errors.Is(err, domain.ErrProjectNotFound) {
				w.WriteHeader(http.StatusOK)
				return
			}
			http.Error(w, "could not resolve project", http.StatusInternalServerError)
			return
		}

		if err := enqueueJob.Execute(r.Context(), project.ID, "process_webhook", string(body), deliveryID); err != nil {
			http.Error(w, "could not enqueue job", http.StatusInternalServerError)
			return
		}

		log.Printf("job enqueued: event=%s delivery=%s project=%s", eventType, deliveryID, project.ID)
		w.WriteHeader(http.StatusOK)
	}
}

func validSignature(body []byte, signature string, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
