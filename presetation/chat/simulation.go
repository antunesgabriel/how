package chat

import (
	"errors"
	"math/rand"
	"time"

	"github.com/antunesgabriel/how-ai/presetation/mock"
	"github.com/antunesgabriel/how-ai/presetation/models"
	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func WaitForAI() tea.Cmd {
	return func() tea.Msg {
		return models.WaitingMsg{}
	}
}

func SimulateAIResponse() tea.Cmd {
	return func() tea.Msg {
		thinkingTime := 1 + rand.Intn(2)
		time.Sleep(time.Duration(thinkingTime) * time.Second)

		responseIndex := rand.Intn(len(mock.AIResponses))
		response := mock.AIResponses[responseIndex]

		return models.AIResponseMsg{
			Content: response,
		}
	}
}

func SimulateTypingResponse(content string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second)

		return models.AIResponseMsg{
			Content: content,
		}
	}
}

func SimulateErrorResponse(errorMsg string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(500 * time.Millisecond)

		var err models.ErrorMsg = errors.New(errorMsg)
		return err
	}
}
