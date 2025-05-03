package tui

import "github.com/charmbracelet/lipgloss"

const (
	ChatModeColor         = "#94E2D5"
	ExecModeColor         = "#FAB387"
	WelcomeBoxBorderColor = "#94E2D5"
	PrimaryTextColor      = "#A6ADC8"
	SecondaryTextColor    = "#6C7086"
)

const (
	ChatIcon = "💬"
	ExecIcon = "🚀"
)

var (
	PromptChatModeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ChatModeColor)).Padding(0, 2)

	PromptExecModeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ExecModeColor)).Padding(0, 2)

	WelcomeBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(WelcomeBoxBorderColor)).
			Padding(1, 2).
			Width(60)

	WelcomeTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ChatModeColor))

	WelcomeSubtitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(SecondaryTextColor))

	WelcomeFooterStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(SecondaryTextColor)).
				MarginTop(1)

	PromptBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ChatModeColor)).
			Padding(1, 2)
)
