package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.String() {

		case tea.KeyEsc.String(), tea.KeyCtrlC.String():
			return m, tea.Quit

		case tea.KeyDown.String(), "j":
			return m, nil

		case tea.KeyUp.String(), "k":
			return m, nil

		case tea.KeyEnter.String():
			return m, nil
		}
	}

	m.prompt, cmd = m.prompt.Update(msg)

	return m, cmd
}
