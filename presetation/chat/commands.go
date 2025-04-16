package chat

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/antunesgabriel/how-ai/presetation/models"
	tea "github.com/charmbracelet/bubbletea"
)

// commandPattern matches Markdown code blocks with bash or shell language tags
// The pattern matches both with and without language specifiers

// commandPattern matches Markdown code blocks with bash or shell language tags
// (?s) makes the dot match newlines as well
var commandPattern = regexp.MustCompile("(?s)```(?:bash|shell)?\\n(.*?)\\n```")

// extractCommandContext tries to find descriptive text around the command
func extractCommandContext(content string, commandMatch string) string {
	// Find the position of the command in the content
	cmdIndex := strings.Index(content, commandMatch)
	if cmdIndex == -1 {
		return ""
	}

	// Extract text before the command (up to 100 characters)
	startIndex := max(0, cmdIndex-100)
	beforeText := content[startIndex:cmdIndex]

	// Find the last paragraph break or header
	lastBreak := -1
	for _, delimiter := range []string{"\n\n", "# ", "## "} {
		if idx := strings.LastIndex(beforeText, delimiter); idx != -1 && idx > lastBreak {
			lastBreak = idx
		}
	}

	if lastBreak != -1 {
		beforeText = beforeText[lastBreak:]
	}

	// Clean up and format
	return strings.TrimSpace(beforeText)
}

// parseCommands extracts commands from the AI response
func parseCommands(content string) []models.Command {
	matches := commandPattern.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	var commands []models.Command
	for _, match := range matches {
		if len(match) > 1 {
			// Extract raw command and trim whitespace
			rawCmd := strings.TrimSpace(match[1])
			if rawCmd == "" {
				continue
			}

			// Split multiple commands if separated by newlines
			cmdLines := strings.Split(rawCmd, "\n")

			for _, cmdLine := range cmdLines {
				// Skip empty lines and comments
				cmdLine = strings.TrimSpace(cmdLine)
				if cmdLine == "" || strings.HasPrefix(cmdLine, "#") {
					continue
				}

				// Get context for the command
				context := extractCommandContext(content, match[0])
				// Create command object
				cmd := models.Command{
					Raw:         cmdLine,
					Description: context,
					Executed:    false,
					Status:      "pending",
					Timestamp:   time.Time{}, // Will be set when executed
					ExitCode:    -1,          // Default until executed
				}
				// Validate the command
				if cmd.ValidateCommand() {
					commands = append(commands, cmd)
				}
			}
		}
	}

	return commands
}

func isSafeCommand(cmd string) bool {
	// Create a temporary Command struct and use its validation method
	tmpCmd := models.Command{
		Raw:    cmd,
		Status: "validating",
	}
	// Use the comprehensive validation from the Command struct
	return tmpCmd.ValidateCommand()
}

// executeCommand runs a shell command and returns the output
func executeCommand(cmd models.Command) tea.Cmd {
	return func() tea.Msg {
		// Add a small delay to show the "Executing command" message
		time.Sleep(300 * time.Millisecond)
		// Validate the command for safety
		if !cmd.ValidateCommand() {
			// Update command status
			cmd.Status = "blocked"
			cmd.Result = "Command blocked for security reasons"
			cmd.ExitCode = -2 // Special code for blocked commands

			return models.CommandResultMsg{
				Command: cmd.Raw,
				Output:  "Command execution blocked for security reasons.",
				Error:   fmt.Errorf("potentially unsafe command"),
			}
		}

		// Update command status
		cmd.Status = "executing"
		cmd.Timestamp = time.Now()
		// Execute the command
		command := exec.Command("sh", "-c", cmd.Raw)
		var stdout, stderr bytes.Buffer
		command.Stdout = &stdout
		command.Stderr = &stderr

		err := command.Run()

		var output string
		if stdout.Len() > 0 {
			output = stdout.String()
		}

		if stderr.Len() > 0 {
			if output != "" {
				output += "\n"
			}
			output += stderr.String()
		}

		// If there's no output but command succeeded, add a message
		if output == "" && err == nil {
			output = "Command executed successfully (no output)"
		}
		// If output is too long, truncate it
		if len(output) > 2000 {
			output = output[:2000] + "\n...\n(Output truncated for display)"
		}
		// Update command status
		if err == nil {
			cmd.Status = "success"
		} else {
			cmd.Status = "failed"
		}
		cmd.Result = output
		cmd.ExitCode = command.ProcessState.ExitCode()

		return models.CommandResultMsg{
			Command: cmd.Raw,
			Output:  output,
			Error:   err,
		}
	}
}

// max returns the larger of x or y
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
