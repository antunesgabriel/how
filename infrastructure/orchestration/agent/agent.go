package einoagent

import (
	einomodel "github.com/cloudwego/eino/components/model"
)

type Agent struct {
	mc einomodel.ToolCallingChatModel
}

func NewAgent(mc einomodel.ToolCallingChatModel) *Agent {
	return &Agent{
		mc: mc,
	}
}
