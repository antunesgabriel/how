package chat

import (
	"fmt"

	"github.com/antunesgabriel/how-ai/presetation/theme"
	"github.com/charmbracelet/lipgloss"
)

func (m Chat) View() string {
	if !m.Ready {
		return "Initializing..."
	}

	if m.Quitting {
		return "Bye!\n"
	}

	var mainView string

	title := theme.TitleStyle.Render("How AI - Terminal Assistant")

	statusBar := ""
	if m.Err != nil {
		statusBar = theme.ErrorStyle.Render(m.Err.Error())
	}

	help := theme.HelpStyle.Render("Enter: Send message • Ctrl+C/Esc: Quit • y/n: Confirm command")

	if m.Waiting {
		spinnerView := fmt.Sprintf("\n  %s Thinking...", m.Spinner.View())
		mainView = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			m.Viewport.View(),
			statusBar,
			spinnerView,
			theme.InputBoxStyle.Width(m.Width-6).Render(m.Textarea.View()),
			help,
		)
	} else {
		mainView = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			m.Viewport.View(),
			statusBar,
			theme.InputBoxStyle.Width(m.Width-6).Render(m.Textarea.View()),
			help,
		)
	}

	view := theme.AppStyle.Width(m.Width).Height(m.Height - 4).Render(mainView)

	if m.ShowConfirmation && m.ActiveCommand != nil {
		return view + "\n" + m.ConfirmDialog.View()
	}

	return view
}

func (m Chat) FullView() string {
	title := theme.TitleStyle.Render("How AI - Terminal Assistant")

	statusBar := ""
	if m.Err != nil {
		statusBar = theme.ErrorStyle.Render(m.Err.Error())
	}

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
