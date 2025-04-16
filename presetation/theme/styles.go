package theme

import "github.com/charmbracelet/lipgloss"

var (
	AppStyle = lipgloss.NewStyle().
			Background(Base).
			Foreground(Text).
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(Lavender).
			Background(Surface0).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1).
			Width(100).
			Align(lipgloss.Center)

	UserMsgStyle = lipgloss.NewStyle().
			Foreground(Green).
			PaddingLeft(1)

	AIMsgStyle = lipgloss.NewStyle().
			Foreground(Mauve).
			PaddingLeft(1)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(Blue)

	InputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Lavender).
			Padding(0, 1).
			MarginTop(1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(Text).
			Background(Surface0).
			Padding(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(Subtext1).
			Italic(true).
			MarginTop(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Red).
			Bold(true).
			MarginTop(1)
)
