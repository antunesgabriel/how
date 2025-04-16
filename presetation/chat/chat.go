package chat

import (
	"time"

	"github.com/antunesgabriel/how-ai/presetation/mock"
	"github.com/antunesgabriel/how-ai/presetation/models"
	"github.com/antunesgabriel/how-ai/presetation/theme"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
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
}

func SimulateAIResponse() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		index := int(t.Unix()) % len(mock.AIResponses)
		return models.AIResponseMsg{Content: mock.AIResponses[index]}
	})
}

func WaitForAI() tea.Cmd {
	return func() tea.Msg {
		return models.WaitingMsg{}
	}
}

func NewChat() (*Chat, error) {
	ta := textarea.New()
	ta.Placeholder = "Type your message here..."
	ta.Focus()
	ta.Prompt = "│ "
	ta.CharLimit = 1000
	ta.SetHeight(3)
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(true)

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

	return &Chat{
		Textarea: ta,
		Messages: messages,
		Viewport: vp,
		Err:      nil,
		Glam:     glam,
		Spinner:  sp,
		Waiting:  false,
	}, nil
}

func (m Chat) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.Spinner.Tick)
}
