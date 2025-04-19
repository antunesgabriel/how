package model

import (
	"context"
	"fmt"

	"github.com/antunesgabriel/how-ai/config"
	einomodel "github.com/cloudwego/eino/components/model"
)

func NewGeminiModel(ctx context.Context, cfg *config.Config) (einomodel.ToolCallingChatModel, error) {
	return nil, fmt.Errorf("gemini model not yet implemented in eino-ext")
}
