package models

// AIResponseMsg is sent when the AI responds to a user query
type AIResponseMsg struct {
	Content string
}

// WaitingMsg is sent when waiting for an AI response
type WaitingMsg struct{}

// ErrorMsg is sent when an error occurs
type ErrorMsg error

// CommandExecutionMsg is sent when a command should be executed
type CommandExecutionMsg struct {
	Command Command
}

// CommandResultMsg is sent when a command has been executed
type CommandResultMsg struct {
	Command string
	Output  string
	Error   error
}

// CommandConfirmationMsg is sent when user confirms/denies command execution
type CommandConfirmationMsg struct {
	Command  Command
	Approved bool
}
