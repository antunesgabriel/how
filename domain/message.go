package domain

type Message struct {
	Role    string
	Content string
}

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)
