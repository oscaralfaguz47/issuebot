package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	httpadapter "github.com/oscaralfaguz47/issuebot/internal/adapter/http"
	"github.com/oscaralfaguz47/issuebot/internal/adapter/postgres"
	"github.com/oscaralfaguz47/issuebot/internal/platform"
	"github.com/oscaralfaguz47/issuebot/internal/usecase"
)

func main() {
	ctx := context.Background()

	// Leo el DSN de la variable de entorno (nunca hardcodeado)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("falta DATABASE_URL")
	}

	// Abro el pool de conexiones
	pool, err := platform.NewDBPool(ctx, dsn)
	if err != nil {
		log.Fatal("no pude conectar a la DB:", err)
	}
	defer pool.Close()

	// Wiring: ACÁ está el cambio. Postgres en vez de memoria.
	repo := postgres.NewProjectRepo(pool)
	createProject := usecase.NewCreateProjectUseCase(repo)

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Post("/projects", httpadapter.HandleCreateProject(createProject))

	log.Println("server en :8080")
	http.ListenAndServe(":8080", r)
}
