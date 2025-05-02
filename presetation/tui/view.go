package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m *model) View() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#eba0ac")).
		Render("** Welcome to How AI **")

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")).
		Render("Use Esc or Ctrl+C to exit")

	return lipgloss.JoinVertical(lipgloss.Left, title, m.prompt.Render(), help)
}
