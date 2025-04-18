package presetation

import (
	"fmt"
	"os"

	"github.com/antunesgabriel/how-ai/infrastructure/openai"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

func StartApp() error {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("this program requires an interactive terminal")
	}

	if os.Getenv("OPENAI_API_KEY") == "" {
		return fmt.Errorf("please set your OPENAI_API_KEY environment variable")
	}

	openAIAgent, err := openai.NewOpenAIAgent()
	if err != nil {
		return fmt.Errorf("failed to create OpenAI agent: %w", err)
	}

	chatModel := NewChatModel(openAIAgent)

	p := tea.NewProgram(chatModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}
