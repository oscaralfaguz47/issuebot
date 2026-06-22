package llm

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicClient struct {
	client anthropic.Client
}

func NewAnthropicClient(apiKey string) *AnthropicClient {
	c := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &AnthropicClient{client: c}
}

func (a *AnthropicClient) GenerateComment(ctx context.Context, issueContent string) (string, error) {
	prompt := "You are a helpful bot. Write a short, friendly comment responding to this GitHub issue:\n\n" + issueContent

	msg, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return "", err
	}

	if len(msg.Content) == 0 {
		return "", nil
	}
	return msg.Content[0].Text, nil
}
