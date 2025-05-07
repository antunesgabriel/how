package tui

import (
	"github.com/antunesgabriel/how/domain"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
)

type model struct {
	agent              domain.Agent
	currentQuery       string
	streaming          bool
	streamingContent   string
	prompt             *prompt
	width              int
	height             int
	glamRenderer       *glamour.TermRenderer
	spinner            spinner.Model
	glamourInitialized bool
	viewport           viewport.Model
	error              error
}

func NewModel(a domain.Agent, initialQuery string) *model {
	p := NewPrompt("Ask me anything...")
	p.Focus()

	s := spinner.New()
	s.Spinner = spinner.MiniDot

	// Initialize viewport with reasonable dimensions
	vp := viewport.New(80, 20)
	vp.SetContent("")

	m := &model{
		agent:              a,
		currentQuery:       initialQuery,
		streamingContent:   "",
		streaming:          false,
		prompt:             p,
		width:              80,
		height:             24,
		glamRenderer:       nil,
		spinner:            s,
		glamourInitialized: false,
		viewport:           vp,
		error:              nil,
	}

	return m
}
