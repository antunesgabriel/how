package presetation

import (
	"fmt"
	"os"

	"github.com/antunesgabriel/how-ai/config"
	"github.com/antunesgabriel/how-ai/domain/agent"
	"github.com/antunesgabriel/how-ai/infrastructure/openai"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

func StartApp(cfg *config.Config, initialQuery string) error {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("this program requires an interactive terminal")
	}

	var llmAgent agent.Agent
	var err error

	switch cfg.DefaultProvider {
	case config.ProviderOpenAI:
		llmAgent, err = openai.NewOpenAIAgentWithConfig(cfg.OpenAI)
	default:
		return fmt.Errorf("unsupported provider: %s", cfg.DefaultProvider)
	}

	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
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
