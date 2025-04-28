package cli

// Message types for internal communication
type (
	// ChatOutputMsg represents a complete AI response
	ChatOutputMsg string

	// PartialOutputMsg represents a partial (streaming) AI response
	PartialOutputMsg string

	// ErrorMsg represents an error message
	ErrorMsg string

	// ViewportContentMsg represents content to be displayed in the viewport
	ViewportContentMsg string
)
