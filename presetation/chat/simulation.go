package chat

import (
	"errors"
	"math/rand"
	"time"

	"github.com/antunesgabriel/how-ai/presetation/mock"
	"github.com/antunesgabriel/how-ai/presetation/models"
	tea "github.com/charmbracelet/bubbletea"
)

// Initialize random seed at package level
func init() {
	// Seed random number generator with current time
	rand.Seed(time.Now().UnixNano())
}

// WaitForAI returns a command that simulates waiting for an AI response
func WaitForAI() tea.Cmd {
	return func() tea.Msg {
		return models.WaitingMsg{}
	}
}

// SimulateAIResponse returns a command that simulates an AI response
func SimulateAIResponse() tea.Cmd {
	return func() tea.Msg {
		// Use variable thinking time
		// Simulate thinking time (1-3 seconds)
		thinkingTime := 1 + rand.Intn(2)
		time.Sleep(time.Duration(thinkingTime) * time.Second)

		// Randomly select a response from the mock responses
		responseIndex := rand.Intn(len(mock.AIResponses))
		response := mock.AIResponses[responseIndex]

		// Return a properly typed message
		return models.AIResponseMsg{
			Content: response,
		}
	}
}

// SimulateTypingResponse returns a command that simulates an AI typing a response character by character
// Note: This isn't used in the current implementation but could be added for a more realistic effect
func SimulateTypingResponse(content string) tea.Cmd {
	return func() tea.Msg {
		// Use variable thinking time
		// Simulate thinking time
		time.Sleep(time.Second)

		// Return a properly typed message
		return models.AIResponseMsg{
			Content: content,
		}
	}
}

// SimulateErrorResponse simulates an error in the AI response
func SimulateErrorResponse(errorMsg string) tea.Cmd {
	return func() tea.Msg {
		// Add a delay to make the error seem more realistic
		time.Sleep(500 * time.Millisecond)

		// Create and return a new error message
		var err models.ErrorMsg = errors.New(errorMsg)
		return err
	}
}
