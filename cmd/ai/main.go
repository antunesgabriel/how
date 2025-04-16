package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/term"
)

var (
	// Styling
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF75B5")).
			MarginLeft(2)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9A9A9A")).
			MarginLeft(2)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			MarginLeft(2)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			MarginLeft(2)

	userStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	assistantStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B5"))

	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B5"))

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Background(lipgloss.Color("#2A2A2A")).
			Padding(0, 1)

	confirmStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF5F00")).
			Padding(0, 1)
)

type model struct {
	textInput      textinput.Model
	messages       []message
	viewport       viewport.Model
	spinner        spinner.Model
	openAIClient   *openai.Client
	renderer       *glamour.TermRenderer
	waitingForAI   bool
	pendingCommand string
	confirmMode    bool
	error          string
	ready          bool
	width          int
	height         int
}

type message struct {
	role    string
	content string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Ask a question or type a command with 'run:' prefix..."
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 80

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	return model{
		textInput:    ti,
		messages:     []message{},
		spinner:      s,
		openAIClient: client,
		renderer:     renderer,
		waitingForAI: false,
		confirmMode:  false,
		error:        "",
		ready:        false,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					return m, executeCommand(m.pendingCommand)
				}

				// Update viewport with the canceled message
				m.messages = append(m.messages, message{
					role:    "system",
					content: "Command execution canceled.",
				})
				return m, updateViewportContent(m)
			}

			input := m.textInput.Value()
			if input == "" {
				return m, nil
			}

			// Add user message to the display
			m.messages = append(m.messages, message{
				role:    "user",
				content: input,
			})

			// Check if it's a command to run
			if strings.HasPrefix(input, "run:") {
				command := strings.TrimSpace(strings.TrimPrefix(input, "run:"))
				m.pendingCommand = command
				m.confirmMode = true
				m.textInput.SetValue("")
				m.textInput.Placeholder = "Execute command? (y/n)"

				// Update viewport with confirmation request
				m.messages = append(m.messages, message{
					role:    "system",
					content: fmt.Sprintf("Do you want to execute: %s", commandStyle.Render(command)),
				})
				return m, updateViewportContent(m)
			}

			m.textInput.SetValue("")
			m.waitingForAI = true

			// Update the viewport
			cmds = append(cmds, updateViewportContent(m))
			cmds = append(cmds, getAIResponse(input))

			return m, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		headerHeight := 6
		footerHeight := 3
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// First time sizing
			m.width = msg.Width
			m.height = msg.Height
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.SetContent("")
			m.ready = true

			// Update input width
			m.textInput.Width = msg.Width - 4
		} else {
			// Resize viewport
			m.width = msg.Width
			m.height = msg.Height
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			m.textInput.Width = msg.Width - 4
		}

		return m, updateViewportContent(m)

	case aiResponseMsg:
		m.waitingForAI = false
		m.messages = append(m.messages, message{
			role:    "assistant",
			content: string(msg),
		})
		return m, updateViewportContent(m)

	case commandOutputMsg:
		m.messages = append(m.messages, message{
			role:    "system",
			content: string(msg),
		})
		return m, updateViewportContent(m)

	case errorMsg:
		m.waitingForAI = false
		m.error = string(msg)
		return m, updateViewportContent(m)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.spinner, spCmd = m.spinner.Update(msg)

	cmds = append(cmds, tiCmd, vpCmd, spCmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	s := titleStyle.Render("Terminal AI Chat") + "\n"
	s += infoStyle.Render("Type your question or use 'run:' prefix to execute a command.") + "\n\n"

	// Display any errors
	if m.error != "" {
		s += errorStyle.Render("Error: "+m.error) + "\n\n"
	}

	// Main content area (viewport)
	s += m.viewport.View() + "\n\n"

	// Input prompt
	promptText := ""
	if m.waitingForAI {
		promptText = m.spinner.View() + " "
	} else if m.confirmMode {
		promptText = confirmStyle.Render("Confirm") + " "
	}

	s += promptStyle.Render(promptText) + m.textInput.View()
	return s
}

// Custom message types
type aiResponseMsg string
type commandOutputMsg string
type errorMsg string

// Get AI response command
func getAIResponse(input string) tea.Cmd {
	return func() tea.Msg {
		// Normally this would use the OpenAI API, but for this example we'll simulate a response
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return errorMsg("OPENAI_API_KEY environment variable not set")
		}

		client := openai.NewClient(apiKey)

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: input,
					},
				},
				MaxTokens: 500,
			},
		)

		if err != nil {
			return errorMsg(fmt.Sprintf("OpenAI API error: %v", err))
		}

		if len(resp.Choices) > 0 {
			return aiResponseMsg(resp.Choices[0].Message.Content)
		}

		return aiResponseMsg("I'm sorry, I couldn't generate a response.")
	}
}

// Execute command with user confirmation
func executeCommand(command string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("sh", "-c", command)
		output, err := cmd.CombinedOutput()

		if err != nil {
			return commandOutputMsg(fmt.Sprintf("Error executing command: %v\n%s", err, output))
		}

		return commandOutputMsg(fmt.Sprintf("Command output:\n%s", output))
	}
}

// Update viewport content
func updateViewportContent(m model) tea.Cmd {
	return func() tea.Msg {
		var content strings.Builder

		for _, msg := range m.messages {
			switch msg.role {
			case "user":
				content.WriteString(userStyle.Render("You: ") + msg.content + "\n\n")
			case "assistant":
				rendered, _ := m.renderer.Render(assistantStyle.Render("AI: ") + msg.content)
				content.WriteString(rendered + "\n\n")
			case "system":
				content.WriteString(msg.content + "\n\n")
			}
		}

		m.viewport.SetContent(content.String())
		m.viewport.GotoBottom()

		return nil
	}
}

func main() {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Println("This program requires an interactive terminal.")
		os.Exit(1)
	}

	// Check for API key
	if os.Getenv("OPENAI_API_KEY") == "" {
		fmt.Println("Please set your OPENAI_API_KEY environment variable.")
		fmt.Println("Example: export OPENAI_API_KEY='your-api-key'")
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
