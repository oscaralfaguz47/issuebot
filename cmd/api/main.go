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

	// Leo el DSN de la variable de entorno (nunca hardcodeado)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Missing DATABASE_URL")
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
	repo := postgres.NewProjectRepo(pool)
	createProject := usecase.NewCreateProjectUseCase(repo)

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Method("POST", "/projects", httpadapter.RequireAuth(jwks, httpadapter.HandleCreateProject(createProject)))

	log.Println("server en :8080")
	http.ListenAndServe(":8080", r)
}
