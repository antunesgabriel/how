package tui

import (
	"github.com/antunesgabriel/how/domain"
)

type model struct {
	agent        domain.Agent
	currentQuery string
	streaming    bool
	prompt       *prompt
	width        int
	height       int
}

func NewModel(a domain.Agent, initialQuery string) *model {
	p := NewPrompt("Ask me anything...")

	p.Focus()

	// Set initial dimensions, but these will be updated by WindowSizeMsg
	return &model{
		agent:        a,
		currentQuery: initialQuery,
		streaming:    false,
		prompt:       p,
		width:        80,
		height:       24,
	}
}
