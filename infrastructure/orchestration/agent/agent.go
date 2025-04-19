package agent

import (
	"context"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type Agent struct {
	runnable compose.Runnable[[]*schema.Message, []*schema.Message]
}

func (a *Agent) GetResponse(ctx context.Context, input string) (string, error) {
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: "You are a helpful assistant.",
		},
		{
			Role:    schema.User,
			Content: input,
		},
	}

	resp, err := a.runnable.Invoke(ctx, messages)
	if err != nil {
		return "I'm sorry, I couldn't generate a response.", err
	}

	if len(resp) == 0 {
		return "I'm sorry, I couldn't generate a response.", nil
	}

	content := ""

	for _, message := range resp {
		content += message.Content
	}

	return content, nil
}

func NewAgent(ctx context.Context, mc einomodel.ToolCallingChatModel) (*Agent, error) {
	searchTool, err := duckduckgo.NewTool(ctx, &duckduckgo.Config{})
	if err != nil {
		return nil, err
	}

	tools := []tool.BaseTool{searchTool}

	toolInfos := make([]*schema.ToolInfo, 0, len(tools))
	for _, tool := range tools {
		info, err := tool.Info(ctx)
		if err != nil {
			return nil, err
		}
		toolInfos = append(toolInfos, info)
	}

	mcWithTools, err := mc.WithTools(toolInfos)
	if err != nil {
		return nil, err
	}

	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: tools,
	})

	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.AppendChatModel(mcWithTools, compose.WithNodeName("chat_model"))
	chain.AppendToolsNode(toolsNode, compose.WithNodeName("tools"))

	runnable, err := chain.Compile(ctx)
	if err != nil {
		return nil, err
	}

	return &Agent{
		runnable: runnable,
	}, nil
}
