package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/antunesgabriel/how-ai/domain/agent"
	openailib "github.com/sashabaranov/go-openai"
)

var _ agent.Agent = (*OpenAIAgent)(nil)

type OpenAIAgent struct {
	client *openailib.Client
}

func NewOpenAIAgent() (*OpenAIAgent, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openailib.NewClient(apiKey)
	return &OpenAIAgent{
		client: client,
	}, nil
}

func (a *OpenAIAgent) GetResponse(ctx context.Context, input string) (string, error) {
	resp, err := a.client.CreateChatCompletion(
		ctx,
		openailib.ChatCompletionRequest{
			Model: openailib.GPT3Dot5Turbo,
			Messages: []openailib.ChatCompletionMessage{
				{
					Role:    openailib.ChatMessageRoleUser,
					Content: input,
				},
			},
			MaxTokens: 500,
		},
	)

	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %v", err)
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "I'm sorry, I couldn't generate a response.", nil
}
