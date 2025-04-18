package model

import (
	"context"
	"time"

	"github.com/antunesgabriel/how-ai/config"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	einomodel "github.com/cloudwego/eino/components/model"
)

func NewOllamaModel(ctx context.Context, cfg *config.Config) (einomodel.ToolCallingChatModel, error) {
	modelCfg := ollama.ChatModelConfig{
		BaseURL: cfg.Ollama.BaseURL,
		Timeout: 30 * time.Second,
		Model:   cfg.Ollama.DefaultModel,
	}

	cm, err := ollama.NewChatModel(ctx, &modelCfg)
	if err != nil {
		return nil, err
	}

	return cm, nil
}
