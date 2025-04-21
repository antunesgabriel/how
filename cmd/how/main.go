package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	einomodel "github.com/cloudwego/eino/components/model"

	"github.com/antunesgabriel/how/config"
	"github.com/antunesgabriel/how/infrastructure/orchestration/agent"
	llmodel "github.com/antunesgabriel/how/infrastructure/orchestration/model"
	"github.com/antunesgabriel/how/presetation"
)

var (
	provider = "" // Provider to use. Exe: openai, claude, gemini, deepseek, ollama
	model    = "" // Provider model to use. Exe: gpt-4o, gpt-3.5-turbo, etc.
)

func main() {
	ctx := context.Background()

	if len(os.Args) > 1 {
		cmd := os.Args[1]

		if cmd == "init" {
			isLocal := slices.Contains(os.Args[2:], "--local")

			if err := handleInit(isLocal); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		query := strings.Join(os.Args[1:], " ")
		if err := startApp(ctx, query); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := startApp(ctx, ""); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func handleInit(isLocal bool) error {
	var configPath string
	var createConfigFunc func() error

	if isLocal {
		configPath = config.LocalConfigFilePath()
		createConfigFunc = config.CreateLocalExampleConfig
	} else {
		configPath = config.GlobalConfigFilePath()
		createConfigFunc = config.CreateGlobalExampleConfig
	}

	_, err := os.Stat(configPath)
	if err == nil {
		fmt.Printf("Configuration file already exists at %s\n", configPath)
		fmt.Println("To create an example configuration with all providers, run: how example")
		return nil
	}

	if err := createConfigFunc(); err != nil {
		return err
	}

	fmt.Printf("Default configuration created at %s\n", configPath)
	fmt.Println("Please edit this file with your API key and preferences.")
	return nil
}

func startApp(ctx context.Context, query string) error {
	cfg, err := config.Load(provider, model)
	if err != nil {
		if strings.Contains(err.Error(), "config file not found") {
			fmt.Println("Configuration file not found.")
			fmt.Println("Run 'how init' to create a default configuration.")
			return fmt.Errorf("configuration required")
		}
		return err
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
		return fmt.Errorf("unsupported provider: %s", cfg.DefaultProvider)
	}

	llmAgent, err := agent.NewAgent(ctx, chatModel)
	if err := presetation.StartApp(llmAgent, query); err != nil {
		return err
	}

	return nil
}
