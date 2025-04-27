package cli

import (
	"github.com/charmbracelet/lipgloss"
)

// Mode represents the CLI mode (Chat or Exec)
type Mode int

const (
	// ChatMode is for asking questions
	ChatMode Mode = iota

	// ExecMode is for executing commands (not implemented yet)
	ExecMode
)

// UI Styles
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

	// PromptStyle is used for the input prompt
	PromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			MarginLeft(2)

	// ErrorStyle is used for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			MarginLeft(2)

	// UserStyle is used for user messages
	UserStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	// AssistantStyle is used for assistant messages
	AssistantStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B5")).
			Bold(true)

	// SystemStyle is used for system messages
	SystemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9A9A9A"))

	// SpinnerStyle is used for the loading spinner
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B5"))

	// ChatModeStyle is used to indicate chat mode
	ChatModeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#04B575")).
			Padding(0, 1)

	// ExecModeStyle is used to indicate exec mode
	ExecModeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF75B5")).
			Padding(0, 1)
)
