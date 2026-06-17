package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	httpadapter "github.com/oscaralfaguz47/issuebot/internal/adapter/http"
	"github.com/oscaralfaguz47/issuebot/internal/adapter/memory"
	"github.com/oscaralfaguz47/issuebot/internal/usecase"
)

func main() {
	// Wiring: armás las dependencias de adentro hacia afuera
	repo := memory.NewProjectRepo()
	createProject := usecase.NewCreateProjectUseCase(repo)

	// Router
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Post("/projects", httpadapter.HandleCreateProject(createProject))

	http.ListenAndServe(":8080", r)
}
