package llm

import (
	"context"
	"fmt"
)

type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) GenerateComment(ctx context.Context, issueContent string) (string, error) {
	return fmt.Sprintf("🤖 [mock] Thanks for opening this issue! (received %d chars)", len(issueContent)), nil
}