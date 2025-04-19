package model

import (
	"context"

	"github.com/antunesgabriel/how-ai/config"
	"github.com/cloudwego/eino-ext/components/model/claude"
	einomodel "github.com/cloudwego/eino/components/model"
)

func NewAnthropicModel(ctx context.Context, cfg *config.Config) (einomodel.ToolCallingChatModel, error) {
	modelCfg := claude.Config{
		APIKey:    cfg.Claude.APIKey,
		Model:     cfg.Claude.Model,
		MaxTokens: cfg.Claude.MaxTokens,

		ByBedrock:       cfg.Claude.ByBedrock,
		AccessKey:       cfg.Claude.AccessKey,
		SecretAccessKey: cfg.Claude.SecretAccessKey,
		SessionToken:    cfg.Claude.SessionToken,
		Region:          cfg.Claude.Region,
	}

	if cfg.Claude.BaseURL != nil {
		modelCfg.BaseURL = cfg.Claude.BaseURL
	}

	model, err := claude.NewChatModel(ctx, &modelCfg)
	if err != nil {
		return nil, err
	}

	return model, nil
}
