package tui

import tea "github.com/charmbracelet/bubbletea"

func (m *model) processStreamContent() tea.Cmd {
	return func() tea.Msg {
		channel := m.agent.GetChannel()
		value := <-channel

		if value.IsLast() {
			return FinishStreamContentMsg{
				Content: value.GetContent(),
			}
		}

		m.streamingContent = value.GetContent()

		return value
	}
}
