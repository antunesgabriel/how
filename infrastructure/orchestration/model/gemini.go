package model

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/gemini"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"github.com/antunesgabriel/how/config"
)

func NewGeminiModel(
	ctx context.Context,
	cfg *config.Config,
) (einomodel.ToolCallingChatModel, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.Gemini.APIKey))
	if err != nil {
		return nil, err
	}

	modelCfg := gemini.Config{
		Client: client,
		Model:  cfg.Gemini.Model,
	}

	if cfg.Gemini.MaxTokens != nil && *cfg.Gemini.MaxTokens > 0 {
		modelCfg.MaxTokens = cfg.Gemini.MaxTokens
	}

	mc, err := gemini.NewChatModel(ctx, &modelCfg)
	if err != nil {
		return nil, err
	}

	return mc, fmt.Errorf("gemini model not yet implemented in eino-ext")
}
