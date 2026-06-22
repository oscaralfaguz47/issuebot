package domain

import "context"

type LLMClient interface {
	GenerateComment(ctx context.Context, issueContent string) (string, error)
}
