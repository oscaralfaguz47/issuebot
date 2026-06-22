package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/oscaralfaguz47/issuebot/internal/adapter/postgres"
	"github.com/oscaralfaguz47/issuebot/internal/domain"
	"github.com/oscaralfaguz47/issuebot/internal/platform"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	pool, err := platform.NewDBPool(ctx, dsn)
	if err != nil {
		log.Fatal("db connection failed:", err)
	}
	defer pool.Close()

	jobRepo := postgres.NewJobRepo(pool)
	log.Println("worker started, polling for jobs...")

	for {
		job, err := jobRepo.ClaimJob(ctx)
		if errors.Is(err, domain.ErrNoJobs) {
			time.Sleep(2 * time.Second) // queue empty, wait and retry
			continue
		}
		if err != nil {
			log.Println("error claiming job:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// --- process the job (placeholder for now) ---
		log.Printf("processing job id=%s type=%s project=%s", job.ID, job.Type, job.ProjectID)

		// mark it done
		if err := jobRepo.MarkDone(ctx, job.ID); err != nil {
			log.Println("error marking done:", err)
			continue
		}
		log.Printf("job done: id=%s", job.ID)
	}
}
