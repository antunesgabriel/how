package agent

import (
	"context"
	"errors"
	"io"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"

	"github.com/antunesgabriel/how/domain"
)

type Agent struct {
	agent   *react.Agent
	history []*schema.Message
	channel chan domain.ChatOutput
	running bool
	mode    domain.AgentMode
}

func NewAgent(
	ctx context.Context,
	toolCallingChatModel einomodel.ToolCallingChatModel,
) (domain.Agent, error) {
	searchOnWeb, err := duckduckgo.NewTool(ctx, &duckduckgo.Config{
		ToolName:   "search_on_web",
		MaxResults: 3,
	})
	if err != nil {
		return nil, err
	}

	toolsConfig := compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{
			searchOnWeb,
		},
	}

	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: toolCallingChatModel,
		ToolsConfig:      toolsConfig,
	})
	if err != nil {
		return nil, err
	}

	return &Agent{
		agent:   agent,
		history: make([]*schema.Message, 0),
		channel: make(chan domain.ChatOutput),
		running: false,
		mode:    domain.ChatAgentMode,
	}, nil
}

func (a *Agent) ChangeMode(mode domain.AgentMode) {
	a.mode = mode
}

func (a *Agent) GetHistory() []domain.Message {
	messages := make([]domain.Message, 0)

	for _, msg := range a.history {
		switch msg.Role {
		case schema.User:
			messages = append(messages, domain.Message{
				Role:    domain.RoleUser,
				Content: msg.Content,
			})
		case schema.Assistant:
			messages = append(messages, domain.Message{
				Role:    domain.RoleAssistant,
				Content: msg.Content,
			})
		case schema.System:
		case schema.Tool:
			continue
		}
	}

	return messages
}

func (a *Agent) GetChannel() chan domain.ChatOutput {
	return a.channel
}

func (a *Agent) Ask(input string) error {
	ctx := context.Background()

	a.appendUserMessage(input)

	var msgReader *schema.StreamReader[*schema.Message]

	a.running = true
	history := a.getCompletionHistory()

	msgReader, err := a.agent.Stream(ctx, history)
	if err != nil {
		a.running = false

		return err
	}
	defer msgReader.Close()

	for {
		msg, err := msgReader.Recv()
		if errors.Is(err, io.EOF) {
			a.channel <- domain.NewChatOutput("", true)

			a.appendAssistantMessage(msg)

			a.running = false

			return nil
		}

		if err != nil {
			a.running = false
			return err
		}

		a.channel <- domain.NewChatOutput(msg.Content, false)

		if !a.running {
			return nil
		}
	}
}

func (a *Agent) appendUserMessage(content string) *Agent {
	a.history = append(a.history, &schema.Message{
		Role:    schema.User,
		Content: content,
	})

	return a
}

func (a *Agent) appendAssistantMessage(msg *schema.Message) *Agent {
	a.history = append(a.history, msg)

	return a
}

func (a *Agent) getCompletionHistory() []*schema.Message {
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: domain.ChatModeSystemPrompt,
		},
	}

	messages = append(messages, a.history...)

	return messages
}
