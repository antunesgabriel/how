package domain

type AgentMode int

const (
	ChatAgentMode AgentMode = iota
)

func (m AgentMode) String() string {
	if m == ChatAgentMode {
		return "chat"
	} else {
		return "unknown"
	}
}

type StreamResponse interface {
	Content() (string, bool, error)
}

type Agent interface {
	Ask(input string) error
	GetChannel() chan ChatOutput
	GetHistory() []Message
	ChangeMode(mode AgentMode)
}

type ChatOutput struct {
	content string
	last    bool
}

func (co ChatOutput) GetContent() string {
	return co.content
}

func (co ChatOutput) IsLast() bool {
	return co.last
}

func NewChatOutput(content string, last bool) ChatOutput {
	return ChatOutput{
		content: content,
		last:    last,
	}
}
