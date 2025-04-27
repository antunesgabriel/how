package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

	"github.com/antunesgabriel/how/domain"
)

// Model represents the main CLI application model
type Model struct {
	// State
	mode          Mode
	waitingForAI  bool
	ready         bool
	error         string
	initialQuery  string
	viewportReady bool
	streamBuffer  string

	// Components
	textInput textinput.Model
	viewport  viewport.Model
	spinner   spinner.Model
	renderer  *glamour.TermRenderer

	// Data
	messages []domain.Message
	agent    domain.Agent
	width    int
	height   int
}

// NewModel creates a new CLI model
func NewModel(agent domain.Agent) *Model {
	ti := textinput.New()
	ti.Placeholder = "Ask me something..."
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
			Content: "Welcome to How CLI! Type a message and press Enter to chat with the AI assistant.",
		},
	}

	return &Model{
		textInput:     ti,
		messages:      initialMessages,
		spinner:       s,
		agent:         agent,
		renderer:      renderer,
		waitingForAI:  false,
		error:         "",
		ready:         false,
		viewportReady: false,
		mode:          ChatMode,
		streamBuffer:  "",
	}
}

// SetInitialQuery sets an initial query to be processed on startup
func (m *Model) SetInitialQuery(query string) {
	m.initialQuery = query
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	cmds := []tea.Cmd{textinput.Blink, m.spinner.Tick}

	if m.initialQuery != "" {
		m.messages = append(m.messages, domain.Message{
			Role:    domain.RoleUser,
			Content: m.initialQuery,
		})

		cmds = append(cmds, m.getAIResponse(m.initialQuery))
		m.waitingForAI = true
	}

	return tea.Batch(cmds...)
}

// Update handles all UI updates
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		case tea.KeyTab:
			if m.mode == ChatMode {
				m.mode = ExecMode
				m.textInput.Placeholder = "Execute something... (not implemented yet)"
			} else {
				m.mode = ChatMode
				m.textInput.Placeholder = "Ask me something..."
			}
			return m, m.updateViewportContent()

		case tea.KeyEnter:
			input := m.textInput.Value()
			if input == "" {
				return m, nil
			}

			m.messages = append(m.messages, domain.Message{
				Role:    domain.RoleUser,
				Content: input,
			})

			m.textInput.SetValue("")
			m.waitingForAI = true
			m.streamBuffer = ""

			cmds = append(cmds, tea.Println(input), m.getAIResponse(input), m.receiveAIResponse())

			return m, tea.Batch(cmds...)

		case tea.KeyCtrlL:
			m.messages = []domain.Message{
				{
					Role:    domain.RoleSystem,
					Content: "Welcome to How CLI! Type a message and press Enter to chat with the AI assistant.",
				},
			}
			return m, m.updateViewportContent()

		case tea.KeyCtrlH:
			helpMsg := domain.Message{
				Role:    domain.RoleSystem,
				Content: "**Help**\n- `tab`: switch between chat and exec modes\n- `ctrl+l`: clear screen\n- `ctrl+h`: show this help\n- `ctrl+c`: quit",
			}
			m.messages = append(m.messages, helpMsg)
			return m, m.updateViewportContent()
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
			m.viewportReady = true

			m.textInput.Width = msg.Width - 4
		} else {
			m.width = msg.Width
			m.height = msg.Height
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			m.textInput.Width = msg.Width - 4
		}

		return m, m.updateViewportContent()

	case ChatOutputMsg:
		m.messages = append(m.messages, domain.Message{
			Role:    domain.RoleAssistant,
			Content: string(msg),
		})
		m.waitingForAI = false
		m.streamBuffer = ""
		return m, m.updateViewportContent()

	case PartialOutputMsg:
		m.streamBuffer += string(msg)
		return m, m.updateStreamContent()
	case domain.ChatOutput:
		if msg.IsLast() {
			m.streamBuffer = ""
			m.waitingForAI = false
			m.messages = append(m.messages, domain.Message{
				Role:    domain.RoleAssistant,
				Content: string(msg.GetContent()),
			})
			return m, tea.Batch(m.updateStreamContent(), m.receiveAIResponse())
		} else {
			return m, m.receiveAIResponse()
		}

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
	if m.viewportReady {
		m.viewport, vpCmd = m.viewport.Update(msg)
	}
	m.spinner, spCmd = m.spinner.Update(msg)

	cmds = append(cmds, tiCmd, vpCmd, spCmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	s := TitleStyle.Render("How CLI - AI Assistant") + "\n"

	modeText := "Mode: " + ChatModeStyle.Render(" Chat ")
	if m.mode == ExecMode {
		modeText = "Mode: " + ExecModeStyle.Render(" Exec ")
	}

	s += InfoStyle.Render("Press Tab to switch modes. Currently in "+modeText) + "\n\n"

	if m.error != "" {
		s += ErrorStyle.Render("Error: "+m.error) + "\n\n"
	}

	if m.viewportReady {
		s += m.viewport.View() + "\n\n"
	}

	promptText := ""
	if m.waitingForAI {
		promptText = m.spinner.View() + " "
	}

	s += PromptStyle.Render(promptText) + m.textInput.View()
	return s
}

func (m *Model) getAIResponse(query string) tea.Cmd {
	return func() tea.Msg {
		err := m.agent.Ask(query)
		if err != nil {
			return ErrorMsg(fmt.Sprintf("Error: %v", err))
		}

		return nil
	}
}

func (m *Model) receiveAIResponse() tea.Cmd {
	return func() tea.Msg {
		output := <-m.agent.GetChannel()
		m.streamBuffer += output.GetContent()
		m.waitingForAI = !output.IsLast()

		return output
	}
}

func (m *Model) updateViewportContent() tea.Cmd {
	return func() tea.Msg {
		var content strings.Builder

		for _, msg := range m.messages {
			switch msg.Role {
			case domain.RoleUser:
				content.WriteString(UserStyle.Render("You: ") + msg.Content + "\n\n")
			case domain.RoleAssistant:
				rendered, err := m.renderer.Render(msg.Content)
				if err != nil {
					content.WriteString(AssistantStyle.Render("How: ") + msg.Content + "\n\n")
				} else {
					content.WriteString(AssistantStyle.Render("How: ") + rendered + "\n\n")
				}
			case domain.RoleSystem:
				rendered, err := m.renderer.Render(msg.Content)
				if err != nil {
					content.WriteString(SystemStyle.Render(msg.Content) + "\n\n")
				} else {
					content.WriteString(SystemStyle.Render(rendered) + "\n\n")
				}
			}
		}

		return ViewportContentMsg(content.String())
	}
}

func (m *Model) updateStreamContent() tea.Cmd {
	return func() tea.Msg {
		var content strings.Builder

		for _, msg := range m.messages {
			switch msg.Role {
			case domain.RoleUser:
				content.WriteString(UserStyle.Render("You: ") + msg.Content + "\n\n")
			case domain.RoleAssistant:
				rendered, err := m.renderer.Render(msg.Content)
				if err != nil {
					content.WriteString(AssistantStyle.Render("How: ") + msg.Content + "\n\n")
				} else {
					content.WriteString(AssistantStyle.Render("How: ") + rendered + "\n\n")
				}
			case domain.RoleSystem:
				rendered, err := m.renderer.Render(msg.Content)
				if err != nil {
					content.WriteString(SystemStyle.Render(msg.Content) + "\n\n")
				} else {
					content.WriteString(SystemStyle.Render(rendered) + "\n\n")
				}
			}
		}

		if m.streamBuffer != "" {
			rendered, err := m.renderer.Render(m.streamBuffer)
			if err != nil {
				content.WriteString(AssistantStyle.Render("How: ") + m.streamBuffer)
			} else {
				content.WriteString(AssistantStyle.Render("How: ") + rendered)
			}
		}

		return ViewportContentMsg(content.String())
	}
}
