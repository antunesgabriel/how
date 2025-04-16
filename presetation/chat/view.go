package chat

import (
	"github.com/antunesgabriel/how-ai/presetation/theme"
	"github.com/charmbracelet/lipgloss"
)

func (m Chat) View() string {
	if !m.Ready {
		return "Initializing..."
	}

	if m.Quitting {
		return "Thanks for using How AI! Goodbye!\n"
	}

	title := theme.TitleStyle.Width(m.Width - 4).Render("How AI")

	statusText := "Press esc or ctrl+c to quit"
	if m.Waiting {
		statusText = m.Spinner.View() + " Thinking..."
	}
	statusBar := theme.StatusBarStyle.Width(m.Width - 4).Render(statusText)

	help := theme.HelpStyle.Render("Enter: Send message • Ctrl+C/Esc: Quit")

	chatContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		m.Viewport.View(),
		statusBar,
		theme.InputBoxStyle.Width(m.Width-4).Render(m.Textarea.View()),
		help,
	)

	return theme.AppStyle.Width(m.Width).Height(m.Height).Render(chatContent)
}
