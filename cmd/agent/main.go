package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	einomodel "github.com/cloudwego/eino/components/model"

	"github.com/antunesgabriel/how/config"
	"github.com/antunesgabriel/how/domain"
	"github.com/antunesgabriel/how/infrastructure/orchestration/agent"
	llmodel "github.com/antunesgabriel/how/infrastructure/orchestration/model"
)

func main() {
	ctx := context.Background()

	a, err := startApp(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// TODO: add debug agent
	stream, err := a.GetStreamResponse(ctx, []domain.Message{{Role: domain.RoleUser, Content: "Tell me about the ls command"}})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for {
		content, finished, err := stream.Content()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Print(content)

		if finished {
			break
		}
	}

}

func startApp(ctx context.Context) (domain.Agent, error) {
	cfg, err := config.Load("", "")
	if err != nil {
		if strings.Contains(err.Error(), "config file not found") {
			fmt.Println("Configuration file not found.")
			fmt.Println("Run 'how init' to create a default configuration.")
			return nil, fmt.Errorf("configuration required")
		}
		return nil, err
	}

	var chatModel einomodel.ToolCallingChatModel

	switch cfg.DefaultProvider {
	case config.ProviderOpenAI:
		chatModel, err = llmodel.NewOpenAIModel(ctx, cfg)
	case config.ProviderGemini:
		chatModel, err = llmodel.NewGeminiModel(ctx, cfg)
	case config.ProviderClaude:
		chatModel, err = llmodel.NewClaudeModel(ctx, cfg)
	case config.ProviderDeepseek:
		chatModel, err = llmodel.NewDeepseekModel(ctx, cfg)
	case config.ProviderOllama:
		chatModel, err = llmodel.NewOllamaModel(ctx, cfg)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.DefaultProvider)
	}

	return agent.NewAgent(ctx, chatModel)
}
