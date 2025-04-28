package cli

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

	"github.com/antunesgabriel/how/domain"
)

// StreamMsg is a message sent when we want to start streaming
type StreamMsg struct {
	Input string
}

// StreamTickMsg is a message sent to check for new content in the stream
type StreamTickMsg struct{}

type UiState struct {
	error       error
	promptMode  PromptMode
	querying    bool
	buffer      string
	initialArgs string
}

type UiComponents struct {
	prompt   *Prompt
	renderer *Renderer
	spinner  *Spinner
}

type Ui struct {
	state      UiState
	components UiComponents
	agent      domain.Agent
	history    []string
	historyPos int
	width      int
	height     int
}

func NewUi(agent domain.Agent, initialArgs string) *Ui {
	rand.Seed(time.Now().UnixNano())

	return &Ui{
		state: UiState{
			error:       nil,
			promptMode:  ChatPromptMode,
			querying:    false,
			buffer:      "",
			initialArgs: initialArgs,
		},
		components: UiComponents{
			prompt: NewPrompt(ChatPromptMode),
			renderer: NewRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(80),
			),
			spinner: NewSpinner(),
		},
		agent:      agent,
		history:    make([]string, 0),
		historyPos: -1,
		width:      80,
		height:     24,
	}
}

func (u *Ui) Init() tea.Cmd {
	cmds := []tea.Cmd{
		tea.ClearScreen,
		textinput.Blink,
		u.components.spinner.Tick,
	}

	// If there's an initial query, process it
	if u.state.initialArgs != "" {
		cmds = append(cmds, func() tea.Msg {
			return StreamMsg{Input: u.state.initialArgs}
		})
	}

	return tea.Batch(cmds...)
}

func (u *Ui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds       []tea.Cmd
		promptCmd  tea.Cmd
		spinnerCmd tea.Cmd
	)

	switch msg := msg.(type) {
	// spinner
	case spinner.TickMsg:
		if u.state.querying {
			u.components.spinner, spinnerCmd = u.components.spinner.Update(msg)
			cmds = append(cmds, spinnerCmd)
		}
	// size
	case tea.WindowSizeMsg:
		u.width = msg.Width
		u.height = msg.Height
		u.components.renderer = NewRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(u.width),
		)
	// keyboard
	case tea.KeyMsg:
		switch msg.Type {
		// quit
		case tea.KeyCtrlC:
			return u, tea.Quit
		// history
		case tea.KeyUp, tea.KeyDown:
			if !u.state.querying && len(u.history) > 0 {
				if msg.Type == tea.KeyUp {
					if u.historyPos < len(u.history)-1 {
						u.historyPos++
					}
				} else {
					if u.historyPos > -1 {
						u.historyPos--
					}
				}

				if u.historyPos >= 0 && u.historyPos < len(u.history) {
					u.components.prompt.SetValue(u.history[u.historyPos])
				} else {
					u.components.prompt.SetValue("")
				}

				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(cmds, promptCmd)
			}
		// switch mode
		case tea.KeyTab:
			if !u.state.querying {
				if u.state.promptMode == ChatPromptMode {
					u.state.promptMode = ExecPromptMode
					u.components.prompt.SetMode(ExecPromptMode)
				} else {
					u.state.promptMode = ChatPromptMode
					u.components.prompt.SetMode(ChatPromptMode)
				}
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(cmds, promptCmd, textinput.Blink)
			}
		// enter
		case tea.KeyEnter:
			if !u.state.querying {
				input := u.components.prompt.GetValue()
				if input != "" {
					// Add to history
					u.history = append([]string{input}, u.history...)
					u.historyPos = -1

					inputPrint := u.components.prompt.AsString()
					u.components.prompt.SetValue("")
					u.components.prompt.Blur()
					u.components.prompt, promptCmd = u.components.prompt.Update(msg)

					if u.state.promptMode == ChatPromptMode {
						cmds = append(
							cmds,
							promptCmd,
							tea.Println(inputPrint),
							func() tea.Msg {
								return StreamMsg{Input: input}
							},
						)
					} else {
						// Exec mode not implemented yet
						cmds = append(
							cmds,
							promptCmd,
							tea.Println(inputPrint),
							tea.Println(u.components.renderer.RenderWarning("Exec mode not implemented yet")),
							textinput.Blink,
						)
					}
				}
			}

		// help
		case tea.KeyCtrlH:
			if !u.state.querying {
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.Println(u.components.renderer.RenderContent(u.components.renderer.RenderHelpMessage())),
					textinput.Blink,
				)
			}

		// clear
		case tea.KeyCtrlL:
			if !u.state.querying {
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.ClearScreen,
					textinput.Blink,
				)
			}

		default:
			u.components.prompt.Focus()
			u.components.prompt, promptCmd = u.components.prompt.Update(msg)
			cmds = append(cmds, promptCmd, textinput.Blink)
		}

	// Stream handling
	case StreamMsg:
		u.state.querying = true
		u.state.buffer = ""

		go func() {
			err := u.agent.Ask(msg.Input)
			if err != nil {
				// We can't directly send a message to the Bubble Tea runtime from a goroutine
				// So we'll just print the error
				fmt.Printf("Error: %v\n", err)
			}
		}()

		return u, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return StreamTickMsg{}
		})

	case StreamTickMsg:
		select {
		case output, ok := <-u.agent.GetChannel():
			if !ok {
				// Channel closed
				u.state.querying = false
				u.components.prompt.Focus()
				return u, textinput.Blink
			}

			if output.IsLast() {
				// Final message
				renderedContent := u.components.renderer.RenderContent(u.state.buffer)
				u.state.querying = false
				u.state.buffer = ""
				u.components.prompt.Focus()
				return u, tea.Sequence(
					tea.Println(renderedContent),
					textinput.Blink,
				)
			}

			// Partial message
			u.state.buffer += output.GetContent()
			return u, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
				return StreamTickMsg{}
			})

		default:
			// No message yet, keep checking
			return u, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
				return StreamTickMsg{}
			})
		}

	// errors
	case ErrorMsg:
		u.state.querying = false
		u.state.error = fmt.Errorf("%s", msg)
		u.components.prompt.Focus()
		return u, tea.Sequence(
			tea.Println(u.components.renderer.RenderError(fmt.Sprintf("[error] %s", u.state.error))),
			textinput.Blink,
		)
	}

	return u, tea.Batch(cmds...)
}

func (u *Ui) View() string {
	if u.state.error != nil {
		return u.components.renderer.RenderError(fmt.Sprintf("[error] %s", u.state.error))
	}

	if !u.state.querying {
		modeText := "Mode: " + ChatModeStyle.Render(" Chat ")
		if u.state.promptMode == ExecPromptMode {
			modeText = "Mode: " + ExecModeStyle.Render(" Exec ")
		}

		return fmt.Sprintf(
			"%s\n%s\n\n%s",
			TitleStyle.Render("How CLI - AI Assistant"),
			InfoStyle.Render("Press Tab to switch modes. Currently in "+modeText),
			u.components.prompt.View(),
		)
	}

	return u.components.spinner.View()
}
