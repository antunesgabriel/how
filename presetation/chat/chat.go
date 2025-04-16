package chat

import (
	"github.com/antunesgabriel/how-ai/presetation/confirm"
	"github.com/antunesgabriel/how-ai/presetation/models"
	"github.com/antunesgabriel/how-ai/presetation/theme"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Chat struct {
	Viewport viewport.Model
	Messages []models.Message
	Textarea textarea.Model
	Err      error
	Glam     *glamour.TermRenderer
	Spinner  spinner.Model
	Waiting  bool
	Quitting bool
	Width    int
	Height   int
	Ready    bool

	ShowConfirmation bool
	ConfirmDialog    confirm.Model
	ActiveCommand    *models.Command
}

func NewChat() (*Chat, error) {
	ta := textarea.New()
	ta.Placeholder = "Type your message here..."
	ta.Focus()
	ta.Prompt = ""
	ta.CharLimit = 1000
	ta.SetHeight(2)
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Blink(true)
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(theme.Overlay1)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(0, 0)
	vp.SetContent("")

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = theme.SpinnerStyle

	glam, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return nil, err
	}

	initialMessage := models.Message{
		Sender:  "AI",
		Content: "# Welcome to How AI!\n\nI'm your AI assistant. How can I help you today?",
		IsAI:    true,
	}

	messages := []models.Message{initialMessage}

	confirmDialog := confirm.New()

	return &Chat{
		Textarea:         ta,
		Messages:         messages,
		Viewport:         vp,
		Err:              nil,
		Glam:             glam,
		Spinner:          sp,
		Waiting:          false,
		ShowConfirmation: false,
		ConfirmDialog:    confirmDialog,
		ActiveCommand:    nil,
	}, nil
}

func (m Chat) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.Spinner.Tick)
}
