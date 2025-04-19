package presetation

import (
	"fmt"
	"os"

	"github.com/antunesgabriel/how/config"
	"github.com/antunesgabriel/how/domain"
	"github.com/antunesgabriel/how/infrastructure/openai"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

func StartApp(cfg *config.Config, initialQuery string) error {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("this program requires an interactive terminal")
	}

	var llmAgent domain.Agent
	var err error

	switch cfg.DefaultProvider {
	case config.ProviderOpenAI:
		llmAgent, err = openai.NewOpenAIAgentWithConfig(cfg.OpenAI)
	case config.ProviderGemini:
		// TODO: Implement Gemini agent when available
		return fmt.Errorf("gemini provider not yet implemented")
	case config.ProviderClaude:
		// TODO: Implement Claude agent when available
		return fmt.Errorf("claude provider not yet implemented")
	case config.ProviderDeepseek:
		// TODO: Implement Deepseek agent when available
		return fmt.Errorf("deepseek provider not yet implemented")
	case config.ProviderOllama:
		// TODO: Implement Ollama agent when available
		return fmt.Errorf("ollama provider not yet implemented")
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
