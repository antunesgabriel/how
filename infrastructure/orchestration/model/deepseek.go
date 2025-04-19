package model

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	einomodel "github.com/cloudwego/eino/components/model"

	"github.com/antunesgabriel/how/config"
)

func NewDeepseekModel(
	ctx context.Context,
	cfg *config.Config,
) (einomodel.ToolCallingChatModel, error) {
	modelCfg := deepseek.ChatModelConfig{
		APIKey:  cfg.Deepseek.APIKey,
		Model:   cfg.Deepseek.Model,
		BaseURL: cfg.Deepseek.BaseURL,
	}

	if cfg.Deepseek.MaxTokens > 0 {
		modelCfg.MaxTokens = cfg.Deepseek.MaxTokens
	}

	mc, err := deepseek.NewChatModel(ctx, &modelCfg)
	if err != nil {
		return nil, err
	}

	return mc, nil
}
