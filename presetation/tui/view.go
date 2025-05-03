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
		MarginBottom(2).
		Render("Use ? for open shortcuts")

	promptBox := PromptBoxStyle.
		BorderForeground(lipgloss.Color(borderColor)).
		Width(contentWidth).
		Render(m.prompt.Render())

	mainContainer := lipgloss.NewStyle().Height(m.height).Width(m.width).Padding(0, 1)

	var mainLayout string
	if m.height > 0 {
		promptAndFooterHeight := 8

		mainLayout = lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.PlaceVertical(
				m.height-promptAndFooterHeight,
				lipgloss.Bottom,
				welcomeContainer,
			),
			promptBox,
			"",
			footer,
		)
	} else {
		mainLayout = lipgloss.JoinVertical(
			lipgloss.Left,
			welcomeContainer,
			promptBox,
			"",
			footer,
		)
	}

	return mainContainer.Render(mainLayout)
}
