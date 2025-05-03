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
	prompt             *prompt
	width              int
	height             int
	glamRenderer       *glamour.TermRenderer
	spinner            spinner.Model
	glamourInitialized bool
	viewport           viewport.Model
}

func NewModel(a domain.Agent, initialQuery string) *model {
	p := NewPrompt("Ask me anything...")
	p.Focus()

	s := spinner.New()
	s.Spinner = spinner.MiniDot

	vp := viewport.New(30, 5)

	m := &model{
		agent:              a,
		currentQuery:       initialQuery,
		streaming:          true,
		prompt:             p,
		width:              80,
		height:             24,
		glamRenderer:       nil,
		spinner:            s,
		glamourInitialized: false,
		viewport:           vp,
	}

	return m
}
