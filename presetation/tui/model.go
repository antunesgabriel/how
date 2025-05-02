package tui

import (
	"github.com/antunesgabriel/how/domain"
)

type model struct {
	agent        domain.Agent
	currentQuery string
	streaming    bool
	prompt       *prompt
}

func NewModel(a domain.Agent, initialQuery string) *model {
	p := NewPrompt("Ask me anything...")

	p.Focus()

	return &model{
		agent:        a,
		currentQuery: initialQuery,
		streaming:    false,
		prompt:       p,
	}
}
