package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	httpadapter "github.com/oscaralfaguz47/issuebot/internal/adapter/http"
	"github.com/oscaralfaguz47/issuebot/internal/adapter/postgres"
	"github.com/oscaralfaguz47/issuebot/internal/platform"
	"github.com/oscaralfaguz47/issuebot/internal/usecase"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Println("The .env file was not found, using environment variables")
	}

	// Read the DSN from the environment variable (never hardcoded)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Missing DATABASE_URL")
	}
	// Read the GitHub webhook secret from environment variable
	webhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Fatal("Missing GITHUB_WEBHOOK_SECRET")
	}
	// Abro el pool de conexiones
	pool, err := platform.NewDBPool(ctx, dsn)
	if err != nil {
		log.Fatal("It was not possible to connect to the DB:", err)
	}
	defer pool.Close()

	jwksURL := os.Getenv("SUPABASE_JWKS_URL")
	if jwksURL == "" {
		log.Fatal("Missing SUPABASE_JWKS_URL")
	}
	jwks, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		log.Fatal("It was not possible to load the JWKS:", err)
	}

	// Wiring: ACÁ está el cambio. Postgres en vez de memoria.
	projectRepo := postgres.NewProjectRepo(pool)
	createProject := usecase.NewCreateProjectUseCase(projectRepo)

	jobRepo := postgres.NewJobRepo(pool)
	enqueueJob := usecase.NewEnqueueJobUseCase(jobRepo)

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	// Protected endpoints require authentication
	r.Method("POST", "/projects", httpadapter.RequireAuth(jwks, httpadapter.HandleCreateProject(createProject)))
	// Webhook endpoint for GitHub events
	r.Post("/webhooks/github", httpadapter.HandleGitHubWebhook(webhookSecret, projectRepo, enqueueJob))

	log.Println("server en :8080")
	http.ListenAndServe(":8080", r)
}
