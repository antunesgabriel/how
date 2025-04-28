package cli

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// TitleStyle is used for the application title
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF75B5")).
			MarginLeft(2)

	// InfoStyle is used for informational text
	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9A9A9A")).
			MarginLeft(2)

	// PromptStyle is used for the prompt
	PromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			MarginLeft(2)

	// ErrorStyle is used for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			MarginLeft(2)

	// UserStyle is used for user messages
	UserStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	// AssistantStyle is used for assistant messages
	AssistantStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B5"))

	// SystemStyle is used for system messages
	SystemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9A9A9A"))

	// SpinnerStyle is used for the spinner
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B5"))

	// ChatModeStyle is used for the chat mode indicator
	ChatModeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#66b3ff")).
			Padding(0, 1)

	// ExecModeStyle is used for the exec mode indicator
	ExecModeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#ffa657")).
			Padding(0, 1)
)
