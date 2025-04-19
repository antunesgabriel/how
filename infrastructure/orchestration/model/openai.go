package model

import (
	"context"
	"time"

	"github.com/antunesgabriel/how-ai/config"
	"github.com/cloudwego/eino-ext/components/model/openai"
	einomodel "github.com/cloudwego/eino/components/model"
)

func NewOpenAIModel(ctx context.Context, cfg *config.Config) (einomodel.ToolCallingChatModel, error) {
	timeout := 30 * time.Second
	if cfg.OpenAI.Timeout > 0 {
		timeout = time.Duration(cfg.OpenAI.Timeout) * time.Millisecond
	}

	modelCfg := openai.ChatModelConfig{
		APIKey:  cfg.OpenAI.APIKey,
		Timeout: timeout,
		Model:   cfg.OpenAI.Model,
	}

	if cfg.OpenAI.BaseURL != "" {
		modelCfg.BaseURL = cfg.OpenAI.BaseURL
	}

	if cfg.OpenAI.ByAzure {
		modelCfg.ByAzure = true
		modelCfg.APIVersion = cfg.OpenAI.APIVersion
	}

	if cfg.OpenAI.MaxTokens != nil {
		modelCfg.MaxTokens = cfg.OpenAI.MaxTokens
	}

	model, err := openai.NewChatModel(ctx, &modelCfg)
	if err != nil {
		return nil, err
	}

	return model, nil
}
