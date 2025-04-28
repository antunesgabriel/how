package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"

	"github.com/antunesgabriel/how/domain"
)

// StartCLI starts the CLI application
func StartCLI(agent domain.Agent, initialQuery string) error {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("this program requires an interactive terminal")
	}

	ui := NewUi(agent, initialQuery)
	p := tea.NewProgram(ui, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}
