package presetation

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF75B5")).
			MarginLeft(2)

	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9A9A9A")).
			MarginLeft(2)

	PromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			MarginLeft(2)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			MarginLeft(2)

	UserStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	AssistantStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B5"))

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B5"))

	CommandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Background(lipgloss.Color("#2A2A2A")).
			Padding(0, 1)

	ConfirmStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF5F00")).
			Padding(0, 1)
)
