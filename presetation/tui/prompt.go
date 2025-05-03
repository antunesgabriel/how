package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PromptLeading string

const (
	ExecPromptLeading PromptLeading = "exec"
	ChatPromptLeading PromptLeading = "chat"
)

type prompt struct {
	input        textinput.Model
	modeFeedback PromptLeading
}

func NewPrompt(placeholder string) *prompt {
	i := textinput.New()
	i.Placeholder = placeholder
	i.Focus()

	return &prompt{
		input:        i,
		modeFeedback: ChatPromptLeading,
	}
}

func (p *prompt) Render() string {
	inputStyle := lipgloss.NewStyle()

	switch p.modeFeedback {
	case ExecPromptLeading:
		inputStyle = PromptExecModeStyle
	case ChatPromptLeading:
		inputStyle = PromptChatModeStyle
	}

	promptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ChatModeColor))
	if p.modeFeedback == ExecPromptLeading {
		promptStyle = promptStyle.Foreground(lipgloss.Color(ExecModeColor))
	}

	return fmt.Sprintf(
		"%s %s",
		promptStyle.Render(">"),
		inputStyle.Render(p.input.Value()),
	)
}

func (p *prompt) UseExecMode() *prompt {
	p.modeFeedback = ExecPromptLeading

	return p
}

func (p *prompt) UseChatMode() *prompt {
	p.modeFeedback = ChatPromptLeading

	return p
}

func (p *prompt) Update(msg tea.Msg) (*prompt, tea.Cmd) {
	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)

	return p, cmd
}

func (p *prompt) Focus() *prompt {
	p.input.Focus()

	return p
}

func (p *prompt) Blur() *prompt {
	p.input.Blur()

	return p
}

func (p *prompt) Value() string {
	return p.input.Value()
}

func (p *prompt) SetValue(value string) *prompt {
	p.input.SetValue(value)

	return p
}

func (p *prompt) ChangePlaceholder(placeholder string) *prompt {
	p.input.Placeholder = placeholder

	return p
}
