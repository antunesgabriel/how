package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m *model) View() string {
	width := m.width
	if width < 80 {
		width = 80
	}

	borderColor := ChatModeColor
	if m.prompt.modeFeedback == ExecPromptLeading {
		borderColor = ExecModeColor
	}

	contentWidth := width - 4

	welcomeTitleText := " Welcome to How AI - Beta, your terminal assistant!"
	if m.prompt.modeFeedback == ExecPromptLeading {
		welcomeTitleText = ExecIcon + welcomeTitleText
	} else {
		welcomeTitleText = ChatIcon + welcomeTitleText
	}

	welcomeTitle := ""

	if m.prompt.modeFeedback == ExecPromptLeading {
		welcomeTitle = WelcomeTitleStyle.Foreground(lipgloss.Color(ExecModeColor)).Render(welcomeTitleText)
	} else {
		welcomeTitle = WelcomeTitleStyle.Foreground(lipgloss.Color(ChatModeColor)).Render(welcomeTitleText)
	}

	welcomeSubtitle := WelcomeSubtitleStyle.Render("- Press Tab to change assistant mode\n\n- Use /config to open config menu\n\n\ncwd: Users/you/projects/todo")

	welcomeBox := WelcomeBoxStyle.
		BorderForeground(lipgloss.Color(borderColor)).
		Width(60).
		MarginBottom(3).
		Render(lipgloss.JoinVertical(lipgloss.Left, welcomeTitle, "", welcomeSubtitle))

	welcomeContainer := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Left).
		Render(welcomeBox)

	footer := WelcomeFooterStyle.
		Width(contentWidth).
		MarginBottom(0).
		Render("Use ? for open shortcuts and Tab for change mode")

	promptLeading := ">"

	if m.streaming {
		promptLeading = m.spinner.View()
	}

	promptBox := PromptBoxStyle.
		BorderForeground(lipgloss.Color(borderColor)).
		Width(contentWidth).
		Render(m.prompt.Render(promptLeading))

	mainContainer := lipgloss.NewStyle().Height(m.height).Width(m.width).Padding(0, 1)

	var mainLayout string
	promptAndFooterHeight := 8

	if m.height == 0 {
		return m.spinner.View()
	}

	errorMessage := ""
	if m.error != nil {
		errorMessage = m.error.Error()
	}

	viewportHeight := m.height - promptAndFooterHeight

	if len(m.agent.GetHistory()) > 0 {
		chat := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ChatModeColor)).
			Width(m.viewport.Width).
			Height(viewportHeight - 2).
			Render(m.RenderHistory())

		mainLayout = lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.PlaceVertical(
				viewportHeight,
				lipgloss.Bottom,
				chat,
			),
			promptBox,
			ErrorMessageStyle.Render(errorMessage),
			footer,
		)
	} else {
		mainLayout = lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.PlaceVertical(
				viewportHeight,
				lipgloss.Bottom,
				welcomeContainer,
			),
			promptBox,
			ErrorMessageStyle.Render(errorMessage),
			footer,
		)
	}

	return mainContainer.Render(mainLayout)
}
