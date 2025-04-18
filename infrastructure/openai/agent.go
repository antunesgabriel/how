package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/antunesgabriel/how-ai/config"
	"github.com/antunesgabriel/how-ai/domain"
	openailib "github.com/sashabaranov/go-openai"
)

var _ domain.Agent = (*OpenAIAgent)(nil)

type OpenAIAgent struct {
	client       *openailib.Client
	defaultModel string
}

func NewOpenAIAgent() (*OpenAIAgent, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openailib.NewClient(apiKey)
	return &OpenAIAgent{
		client:       client,
		defaultModel: openailib.GPT3Dot5Turbo,
	}, nil
}

func NewOpenAIAgentWithConfig(cfg *config.OpenAIConfig) (*OpenAIAgent, error) {
	if cfg == nil {
		return nil, fmt.Errorf("OpenAI configuration is required")
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	clientConfig := openailib.DefaultConfig(cfg.APIKey)

	if cfg.BaseURL != "" {
		clientConfig.BaseURL = cfg.BaseURL
	}

	client := openailib.NewClientWithConfig(clientConfig)

	return &OpenAIAgent{
		client:       client,
		defaultModel: cfg.DefaultModel,
	}, nil
}

func (a *OpenAIAgent) GetResponse(ctx context.Context, input string) (string, error) {
	resp, err := a.client.CreateChatCompletion(
		ctx,
		openailib.ChatCompletionRequest{
			Model: a.defaultModel,
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
