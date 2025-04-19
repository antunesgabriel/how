package model

import (
	"context"
	"time"

	"github.com/antunesgabriel/how/config"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	einomodel "github.com/cloudwego/eino/components/model"
)

func NewOllamaModel(ctx context.Context, cfg *config.Config) (einomodel.ToolCallingChatModel, error) {
	timeout := 30 * time.Second
	if cfg.Ollama.Timeout > 0 {
		timeout = time.Duration(cfg.Ollama.Timeout) * time.Millisecond
	}

	modelCfg := ollama.ChatModelConfig{
		BaseURL: cfg.Ollama.BaseURL,
		Timeout: timeout,
		Model:   cfg.Ollama.Model,
	}

	cm, err := ollama.NewChatModel(ctx, &modelCfg)
	if err != nil {
		return nil, err
	}

	return cm, nil
}
