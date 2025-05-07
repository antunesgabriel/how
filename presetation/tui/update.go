package tui

import (
	"time"

	"github.com/antunesgabriel/how/domain"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds       []tea.Cmd
		promptCmd  tea.Cmd
		spinnerCmd tea.Cmd
	)

	switch message := msg.(type) {
	case spinner.TickMsg:
		if m.streaming {
			m.spinner, spinnerCmd = m.spinner.Update(msg)
			cmds = append(cmds, spinnerCmd)
		}
	case tea.WindowSizeMsg:
		m.width = message.Width
		m.height = message.Height
		m.viewport.Width = message.Width - 4
		m.viewport.Height = message.Height - 8

		history := m.agent.GetHistory()
		if len(history) > 0 {
			content := m.RenderHistory()
			m.viewport.SetContent(lipgloss.NewStyle().
				Width(m.viewport.Width).
				Render(content),
			)

			m.viewport.GotoBottom()
		}

	case domain.ChatOutput:
		content := m.RenderHistory()
		m.viewport.SetContent(lipgloss.NewStyle().
			Width(m.viewport.Width).
			Render(content),
		)

		m.viewport.GotoBottom()
		spinnerCmd = m.spinner.Tick
		cmds = append(cmds, spinnerCmd, m.processStreamContent())

		return m, tea.Batch(cmds...)

	case FinishStreamContentMsg:
		m.streaming = false
		m.streamingContent = ""
		m.prompt.Focus()
		m.prompt, promptCmd = m.prompt.Update(msg)

		content := m.RenderHistory()
		m.viewport.SetContent(lipgloss.NewStyle().
			Width(m.viewport.Width).
			Render(content),
		)

		cmds = append(cmds, promptCmd, textinput.Blink)

	case tea.KeyMsg:
		switch message.String() {

		case tea.KeyEsc.String(), tea.KeyCtrlC.String():
			return m, tea.Quit

		case tea.KeyTab.String():
			if m.prompt.modeFeedback == ChatPromptLeading {
				m.prompt.UseExecMode()
			} else {
				m.prompt.UseChatMode()
			}
		case tea.KeyEnter.String():
			if m.streaming {
				return m, nil
			}

			input := m.prompt.Value()
			if input == "" {
				return m, nil
			}

			m.prompt.SetValue("")
			m.prompt.Blur()
			m.prompt, promptCmd = m.prompt.Update(msg)

			go func() {
				err := m.agent.Ask(input)
				if err != nil {
					m.error = err
				}
			}()

			time.Sleep(200 * time.Millisecond)

			content := m.RenderHistory()
			m.viewport.SetContent(lipgloss.NewStyle().
				Width(m.viewport.Width).
				Render(content),
			)

			m.viewport.GotoBottom()

			m.streaming = true
			spinnerCmd = m.spinner.Tick
			cmds = append(cmds, promptCmd, spinnerCmd, textinput.Blink)

			return m, tea.Batch(cmds...)
		case tea.KeyDown.String(), "j":
		case tea.KeyUp.String(), "k":
		default:
			m.prompt.Focus()
			m.prompt, promptCmd = m.prompt.Update(msg)
			cmds = append(cmds, promptCmd, textinput.Blink)
		}

	}

	return m, tea.Batch(cmds...)
}
