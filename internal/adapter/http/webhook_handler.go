package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
)

func HandleGitHubWebhook(webhookSecret string) http.HandlerFunc {
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

		log.Printf("webhook recibido: event=%s delivery=%s body=%d bytes", eventType, deliveryID, len(body))
		w.WriteHeader(http.StatusOK)
	}
}

func validSignature(body []byte, signature string, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
