package agent

import (
	"context"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"

	"github.com/antunesgabriel/how/domain"
)

type Agent struct {
	agent *react.Agent
}

func (a *Agent) GetResponse(ctx context.Context, messages []domain.Message) (string, error) {
	msgs := make([]*schema.Message, len(messages))

	for idx, msg := range messages {
		switch msg.Role {
		case domain.RoleUser:
			msgs[idx] = &schema.Message{
				Role:    schema.User,
				Content: msg.Content,
			}
		case domain.RoleAssistant:
			msgs[idx] = &schema.Message{
				Role:    schema.Assistant,
				Content: msg.Content,
			}
		default:
			msgs[idx] = &schema.Message{
				Role:    schema.System,
				Content: msg.Content,
			}
		}
	}

	outMessage, err := a.agent.Generate(ctx, msgs)
	if err != nil {
		return "Sorry, I get an error when I try to do that", err
	}

	if outMessage == nil {
		return "Sorry, I dont know how to do that", nil
	}

	return outMessage.Content, nil
}

func (a *Agent) GetStreamResponse(
	ctx context.Context,
	messages []domain.Message,
) (domain.StreamResponse, error) {
	msgs := make([]*schema.Message, len(messages))

	for idx, msg := range messages {
		switch msg.Role {
		case domain.RoleUser:
			msgs[idx] = &schema.Message{
				Role:    schema.User,
				Content: msg.Content,
			}
		case domain.RoleAssistant:
			msgs[idx] = &schema.Message{
				Role:    schema.Assistant,
				Content: msg.Content,
			}
		default:
			msgs[idx] = &schema.Message{
				Role:    schema.System,
				Content: msg.Content,
			}
		}
	}

	var msgReader *schema.StreamReader[*schema.Message]

	msgReader, err := a.agent.Stream(ctx, msgs)
	if err != nil {
		return nil, err
	}

	return NewAgentStreamResponse(msgReader), nil
}

func NewAgent(
	ctx context.Context,
	toolCallingChatModel einomodel.ToolCallingChatModel,
) (*Agent, error) {
	searchTool, err := duckduckgo.NewTool(ctx, &duckduckgo.Config{})
	if err != nil {
		return nil, err
	}

	toolsConfig := compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{
			searchTool,
		},
	}

	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: toolCallingChatModel,
		ToolsConfig:      toolsConfig,
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message {
			res := make([]*schema.Message, 0, len(input)+1)

			res = append(
				res,
				schema.SystemMessage(
					"You are an expert in shell commands and terminal operations. Your task is search and to provide detailed, accurate explanations of shell commands that users are considering executing. Break down each part of the command, explain what it does, identify any potential risks or side effects, and explain why someone might want to run it. Be specific about what files or systems will be affected. If the command could potentially be harmful, make sure to clearly highlight those risks.",
				),
			)
			res = append(res, input...)
			return res
		},
	})
	if err != nil {
		return nil, err
	}

	return &Agent{agent: agent}, nil
}
