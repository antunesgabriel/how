package model

import (
	"context"
	"time"

	"github.com/antunesgabriel/how-ai/config"
	"github.com/cloudwego/eino-ext/components/model/openai"
	einomodel "github.com/cloudwego/eino/components/model"
)

func NewOpenAIModel(ctx context.Context, cfg *config.Config) (einomodel.ToolCallingChatModel, error) {
	modelCfg := openai.ChatModelConfig{
		APIKey:  cfg.OpenAI.APIKey,
		Timeout: 30 * time.Second,
		Model:   cfg.OpenAI.DefaultModel,
	}

	if cfg.OpenAI.BaseURL != "" {
		modelCfg.BaseURL = cfg.OpenAI.BaseURL
	}

	model, err := openai.NewChatModel(ctx, &modelCfg)
	if err != nil {
		return nil, err
	}

	return model, nil
}
