package presetation

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"

	"github.com/antunesgabriel/how/domain"
)

func StartApp(llmAgent domain.Agent, initialQuery string) error {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("this program requires an interactive terminal")
	}

	chatModel := NewChatModel(llmAgent)

	if initialQuery != "" {
		chatModel.SetInitialQuery(initialQuery)
	}

	p := tea.NewProgram(chatModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}
