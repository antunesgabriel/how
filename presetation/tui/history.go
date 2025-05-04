package tui

import (
	"fmt"

	"github.com/antunesgabriel/how/domain"
)

func (m *model) RenderHistory() string {
	content := ""
	history := m.agent.GetHistory()

	if len(history) == 0 {
		return content
	}

	for _, msg := range history {
		if msg.Role == domain.RoleUser {
			content += fmt.Sprintf("%s:\n%s\n", UserMessageStyle.Render("You:"), msg.Content)
		} else {
			content += fmt.Sprintf("%s:\n%s\n", AssistantMessageStyle.Render("AI:"), msg.Content)
		}
	}

	if m.streaming {
		content += fmt.Sprintf("%s:\n%s\n", AssistantMessageStyle.Render("AI:"), m.streamingContent)
	}

	return content
}
