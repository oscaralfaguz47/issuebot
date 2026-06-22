package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/oscaralfaguz47/issuebot/internal/adapter/github"
	"github.com/oscaralfaguz47/issuebot/internal/adapter/llm"
	"github.com/oscaralfaguz47/issuebot/internal/adapter/postgres"
	"github.com/oscaralfaguz47/issuebot/internal/domain"
	"github.com/oscaralfaguz47/issuebot/internal/platform"
)

type issuePayload struct {
	Issue struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	} `json:"issue"`
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
	Installation struct {
		ID int64 `json:"id"`
	} `json:"installation"`
}

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
	var llmClient domain.LLMClient
	if os.Getenv("LLM_MODE") == "real" {
		llmClient = llm.NewAnthropicClient(os.Getenv("ANTHROPIC_API_KEY"))
		log.Println("using real Anthropic LLM")
	} else {
		llmClient = llm.NewMockClient()
		log.Println("using mock LLM")
	}

	appID, _ := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	ghClient := github.NewClient(appID, os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH"))

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

		// --- process the job ---
		// parse the payload to get issue/repo/installation info
		var p issuePayload
		if err := json.Unmarshal([]byte(job.Payload), &p); err != nil {
			log.Printf("job %s: bad payload (terminal): %v", job.ID, err)
			jobRepo.MarkFailed(ctx, job.ID) // malformed payload = terminal, no retry
			continue
		}

		// generate the comment from the actual issue title + body
		issueText := p.Issue.Title + "\n\n" + p.Issue.Body
		comment, err := llmClient.GenerateComment(ctx, issueText)
		if err != nil {
			log.Printf("job %s: LLM error (retriable): %v", job.ID, err)
			jobRepo.Reschedule(ctx, job.ID)
			continue
		}

		// post it to GitHub
		err = ghClient.PostComment(ctx, p.Installation.ID, p.Repository.Owner.Login, p.Repository.Name, p.Issue.Number, comment)
		if err != nil {
			handleRetriable(ctx, jobRepo, job, "GitHub post error")
			continue
		}

		if err := jobRepo.MarkDone(ctx, job.ID); err != nil {
			log.Println("error marking done:", err)
			continue
		}
		log.Printf("job done: id=%s, commented on %s/%s#%d", job.ID, p.Repository.Owner.Login, p.Repository.Name, p.Issue.Number)
	}
}

const maxAttempts = 3

func handleRetriable(ctx context.Context, jobRepo *postgres.JobRepo, job *domain.Job, reason string) {
	if job.Attempts+1 >= maxAttempts {
		log.Printf("job %s: %s — giving up after %d attempts (terminal)", job.ID, reason, job.Attempts+1)
		jobRepo.MarkFailed(ctx, job.ID)
		return
	}
	log.Printf("job %s: %s — rescheduling (attempt %d)", job.ID, reason, job.Attempts+1)
	jobRepo.Reschedule(ctx, job.ID)
}
