package presetation

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

	"github.com/antunesgabriel/how/domain"
)

type ChatModel struct {
	textInput      textinput.Model
	messages       []domain.Message
	viewport       viewport.Model
	spinner        spinner.Model
	agent          domain.Agent
	renderer       *glamour.TermRenderer
	waitingForAI   bool
	pendingCommand string
	confirmMode    bool
	error          string
	ready          bool
	width          int
	height         int
	initialQuery   string
}

func NewChatModel(agent domain.Agent) *ChatModel {
	ti := textinput.New()
	ti.Placeholder = "Ask a question or type a command with 'run:' prefix..."
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 80

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	initialMessages := []domain.Message{
		{
			Role:    domain.RoleSystem,
			Content: "Welcome to Terminal AI Chat! Type a message and press Enter to chat with the AI.",
		},
	}

	return &ChatModel{
		textInput:    ti,
		messages:     initialMessages,
		spinner:      s,
		agent:        agent,
		renderer:     renderer,
		waitingForAI: false,
		confirmMode:  false,
		error:        "",
		ready:        false,
	}
}

func (m *ChatModel) SetInitialQuery(query string) {
	m.initialQuery = query
}

func (m *ChatModel) Init() tea.Cmd {
	cmds := []tea.Cmd{textinput.Blink, m.spinner.Tick, m.updateViewportContent()}

	if m.initialQuery != "" {
		m.messages = append(m.messages, domain.Message{
			Role:    domain.RoleUser,
			Content: m.initialQuery,
		})

		cmds = append(cmds, m.getAIResponse())
		m.waitingForAI = true
	}

	return tea.Batch(cmds...)
}

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		spCmd tea.Cmd
		cmds  []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.confirmMode {
				input := strings.ToLower(m.textInput.Value())
				m.confirmMode = false
				m.textInput.SetValue("")
				m.textInput.Placeholder = "Ask a question or type a command with 'run:' prefix..."

				if input == "y" || input == "yes" {
					return m, m.executeCommand(m.pendingCommand)
				}

				m.messages = append(m.messages, domain.Message{
					Role:    domain.RoleSystem,
					Content: "Command execution canceled.",
				})
				return m, m.updateViewportContent()
			}

			input := m.textInput.Value()
			if input == "" {
				return m, nil
			}

			m.messages = append(m.messages, domain.Message{
				Role:    domain.RoleUser,
				Content: input,
			})

			if strings.HasPrefix(input, "run:") {
				command := strings.TrimSpace(strings.TrimPrefix(input, "run:"))
				m.pendingCommand = command
				m.confirmMode = true
				m.textInput.SetValue("")
				m.textInput.Placeholder = "Execute command? (y/n)"

				m.messages = append(m.messages, domain.Message{
					Role:    domain.RoleSystem,
					Content: fmt.Sprintf("Do you want to execute: %s", CommandStyle.Render(command)),
				})
				return m, m.updateViewportContent()
			}

			m.textInput.SetValue("")
			m.waitingForAI = true

			cmds = append(cmds, m.updateViewportContent())
			cmds = append(cmds, m.getAIResponse())

			return m, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		headerHeight := 6
		footerHeight := 3
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.width = msg.Width
			m.height = msg.Height
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.SetContent("")
			m.ready = true

			m.textInput.Width = msg.Width - 4
		} else {
			m.width = msg.Width
			m.height = msg.Height
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			m.textInput.Width = msg.Width - 4
		}

		return m, m.updateViewportContent()

	case AIResponseMsg:
		m.waitingForAI = false
		m.messages = append(m.messages, domain.Message{
			Role:    domain.RoleAssistant,
			Content: string(msg),
		})
		return m, m.updateViewportContent()

	case CommandOutputMsg:
		m.messages = append(m.messages, domain.Message{
			Role:    domain.RoleSystem,
			Content: string(msg),
		})
		return m, m.updateViewportContent()

	case ErrorMsg:
		m.waitingForAI = false
		m.error = string(msg)
		return m, m.updateViewportContent()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case ViewportContentMsg:
		m.viewport.SetContent(string(msg))
		m.viewport.GotoBottom()
		return m, nil
	}

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.spinner, spCmd = m.spinner.Update(msg)

	cmds = append(cmds, tiCmd, vpCmd, spCmd)
	return m, tea.Batch(cmds...)
}

func (m *ChatModel) View() string {
	if !m.ready {
		return "Initializing..."
	}

	s := TitleStyle.Render("How - Terminal AI Assistant") + "\n"
	s += InfoStyle.Render("Type your question to get an answer or request help to How assistant") + "\n\n"

	if m.error != "" {
		s += ErrorStyle.Render("Error: "+m.error) + "\n\n"
	}

	s += m.viewport.View() + "\n\n"

	promptText := ""
	if m.waitingForAI {
		promptText = m.spinner.View() + " "
	} else if m.confirmMode {
		promptText = ConfirmStyle.Render("Confirm") + " "
	}

	s += PromptStyle.Render(promptText) + m.textInput.View()
	return s
}

func (m *ChatModel) getAIResponse() tea.Cmd {
	return func() tea.Msg {
		err := m.agent.Ask("")
		if err != nil {
			return ErrorMsg(fmt.Sprintf("Error: %v", err))
		}

		return nil
	}
}

func (m *ChatModel) executeCommand(command string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("sh", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return CommandOutputMsg(fmt.Sprintf("Error executing command: %v\n%s", err, output))
		}

		return CommandOutputMsg(fmt.Sprintf("Command output:\n%s", output))
	}
}

func (m *ChatModel) updateViewportContent() tea.Cmd {
	return func() tea.Msg {
		var content strings.Builder

		for _, msg := range m.messages {
			switch msg.Role {
			case domain.RoleUser:
				content.WriteString(UserStyle.Render("You: ") + msg.Content + "\n")
			case domain.RoleAssistant:
				rendered, err := m.renderer.Render(msg.Content)
				if err != nil {
					content.WriteString(AssistantStyle.Render("How: ") + msg.Content + "\n")
				}

				content.WriteString(AssistantStyle.Render("How: ") + rendered + "\n")
			case domain.RoleSystem:
				content.WriteString(msg.Content + "\n\n")
			}
		}

		return ViewportContentMsg(content.String())
	}
}
