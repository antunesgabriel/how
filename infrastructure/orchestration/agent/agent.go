package agent

import (
	"context"

	einomodel "github.com/cloudwego/eino/components/model"
)

type Agent struct {
	mc einomodel.ToolCallingChatModel
}

func (a *Agent) GetResponse(ctx context.Context, input string) (string, error) {
	return "I'm sorry, I couldn't generate a response.", nil
}

func NewAgent(mc einomodel.ToolCallingChatModel) (*Agent, error) {
	return &Agent{
		mc: mc,
	}, nil
}
