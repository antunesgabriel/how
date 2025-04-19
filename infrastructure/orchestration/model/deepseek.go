package model

import (
	"context"
	"fmt"

	"github.com/antunesgabriel/how/config"
	einomodel "github.com/cloudwego/eino/components/model"
)

func NewDeepseekModel(ctx context.Context, cfg *config.Config) (einomodel.ToolCallingChatModel, error) {
	return nil, fmt.Errorf("deepseek model not yet implemented in eino-ext")
}
