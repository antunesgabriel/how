package chat

import (
	"strings"

	"github.com/antunesgabriel/how-ai/presetation/input"
	"github.com/antunesgabriel/how-ai/presetation/models"
	"github.com/antunesgabriel/how-ai/presetation/theme"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		spCmd tea.Cmd
		cfCmd tea.Cmd
	)

	if m.ShowConfirmation && m.ActiveCommand != nil {
		m.ConfirmDialog.Width = m.Width
		m.ConfirmDialog.Height = m.Height
		m.ConfirmDialog, cfCmd = m.ConfirmDialog.Update(msg)
		return m, cfCmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		if !m.Ready {
			m.Ready = true
			m.Viewport.Width = msg.Width - 4
			m.Viewport.Height = msg.Height - m.Textarea.Height() - 6
			m.Textarea.SetWidth(msg.Width - 4)
		} else {
			m.Viewport.Width = msg.Width - 4
			m.Viewport.Height = msg.Height - m.Textarea.Height() - 6
			m.Textarea.SetWidth(msg.Width - 4)
		}

		m.updateViewportContent()
		m.Viewport.GotoBottom()

		return m, nil

	case tea.KeyMsg:
		if m.Waiting {
			break
		}

		switch {
		case key.Matches(msg, input.Keymap.Quit):
			m.Quitting = true
			return m, tea.Quit

		case key.Matches(msg, input.Keymap.Send):
			if strings.TrimSpace(m.Textarea.Value()) != "" {
				userMsg := m.Textarea.Value()

				userMessage := models.Message{
					Sender:  "You",
					Content: userMsg,
					IsAI:    false,
				}
				m.Messages = append(m.Messages, userMessage)

				m.Textarea.Reset()
				m.updateViewportContent()
				m.Viewport.GotoBottom()

				m.Waiting = true

				return m, tea.Batch(
					WaitForAI(),
					SimulateAIResponse(),
				)
			}
		}

	case models.AIResponseMsg:
		m.Waiting = false

		rendered, err := m.Glam.Render(msg.Content)
		if err != nil {
			rendered = msg.Content
		}

		commands := parseCommands(msg.Content)

		aiMessage := models.Message{
			Sender:   "AI",
			Content:  rendered,
			IsAI:     true,
			Commands: commands,
		}
		m.Messages = append(m.Messages, aiMessage)

		m.updateViewportContent()
		m.Viewport.GotoBottom()

		// If there are commands, prompt for confirmation of the first one
		if len(commands) > 0 {
			m.ShowConfirmation = true
			m.ActiveCommand = &commands[0]
			m.ConfirmDialog.SetCommand(commands[0])
			return m, nil
		}

		return m, nil

	case models.ErrorMsg:
		m.Err = msg
		return m, nil

	case models.CommandConfirmationMsg:
		m.ShowConfirmation = false

		if msg.Approved {
			executingMsg := models.Message{
				Sender:  "System",
				Content: "Executing command: `" + msg.Command.Raw + "`",
				IsAI:    false,
			}
			m.Messages = append(m.Messages, executingMsg)
			m.updateViewportContent()
			m.Viewport.GotoBottom()

			// Execute the command
			return m, executeCommand(msg.Command)
		} else {
			// Command was rejected
			rejectedMsg := models.Message{
				Sender:  "System",
				Content: "Command execution cancelled.",
				IsAI:    false,
			}
			m.Messages = append(m.Messages, rejectedMsg)
			m.updateViewportContent()
			m.Viewport.GotoBottom()
		}

		m.ActiveCommand = nil
		return m, nil

	case models.CommandResultMsg:
		var resultContent string
		var status string
		var exitCode int

		if msg.Error != nil {
			status = "failed"
			exitCode = 1
			resultContent = models.FormatErrorMessage(msg.Command, msg.Error.Error(), msg.Output)
		} else {
			status = "success"
			exitCode = 0
			resultContent = models.FormatSuccessMessage(msg.Command, msg.Output)
		}

		tmpCmd := models.Command{
			Raw:      msg.Command,
			Status:   status,
			ExitCode: exitCode,
		}

		resultContent += "\n\nStatus: " + tmpCmd.FormattedStatus()
		rendered, err := m.Glam.Render(resultContent)
		if err != nil {
			rendered = resultContent
		}

		resultMsg := models.Message{
			Sender:  "System",
			Content: rendered,
			IsAI:    false,
		}
		m.Messages = append(m.Messages, resultMsg)
		m.updateViewportContent()
		m.Viewport.GotoBottom()

		return m, nil
	}

	m.Textarea, tiCmd = m.Textarea.Update(msg)
	m.Viewport, vpCmd = m.Viewport.Update(msg)
	m.Spinner, spCmd = m.Spinner.Update(msg)

	return m, tea.Batch(tiCmd, vpCmd, spCmd)
}

func (m *Chat) updateViewportContent() {
	var sb strings.Builder

	for i, msg := range m.Messages {
		if i > 0 {
			sb.WriteString("\n\n")
		}

		if msg.IsAI {
			sb.WriteString(theme.AIMsgStyle.Render("AI:"))
			sb.WriteString("\n")
			sb.WriteString(msg.Content)
		} else {
			sb.WriteString(theme.UserMsgStyle.Render("You:"))
			sb.WriteString("\n")
			sb.WriteString(msg.Content)
		}
	}

	m.Viewport.SetContent(sb.String())
}
