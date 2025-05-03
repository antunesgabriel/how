package tui

import "github.com/charmbracelet/lipgloss"

const (
	ChatModeColor         = "#94E2D5"
	ExecModeColor         = "#FAB387"
	WelcomeBoxBorderColor = "#94E2D5" // #94E2D5 if is chat mode and #FAB387 to exec mode
	PrimaryTextColor      = "#A6ADC8"
	SecondaryTextColor    = "#6C7086"
)

var (
	PromptChatModeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ChatModeColor)).Padding(0, 2).MaxWidth(20)

	PromptExecModeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ExecModeColor)).Padding(0, 2).MaxWidth(20)

	WelcomeBoxStyle = lipgloss.NewStyle().Border(lipgloss.BlockBorder())

	WelcomeTitleStyle = lipgloss.NewStyle()
)
