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

var commandPattern = regexp.MustCompile("(?s)```(?:bash|shell)?\\n(.*?)\\n```")

func extractCommandContext(content string, commandMatch string) string {
	cmdIndex := strings.Index(content, commandMatch)
	if cmdIndex == -1 {
		return ""
	}

	startIndex := max(0, cmdIndex-100)
	beforeText := content[startIndex:cmdIndex]

	lastBreak := -1
	for _, delimiter := range []string{"\n\n", "# ", "## "} {
		if idx := strings.LastIndex(beforeText, delimiter); idx != -1 && idx > lastBreak {
			lastBreak = idx
		}
	}

	if lastBreak != -1 {
		beforeText = beforeText[lastBreak:]
	}

	return strings.TrimSpace(beforeText)
}

func parseCommands(content string) []models.Command {
	matches := commandPattern.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	var commands []models.Command
	for _, match := range matches {
		if len(match) > 1 {
			rawCmd := strings.TrimSpace(match[1])
			if rawCmd == "" {
				continue
			}

			cmdLines := strings.Split(rawCmd, "\n")

			for _, cmdLine := range cmdLines {
				cmdLine = strings.TrimSpace(cmdLine)
				if cmdLine == "" || strings.HasPrefix(cmdLine, "#") {
					continue
				}

				context := extractCommandContext(content, match[0])
				cmd := models.Command{
					Raw:         cmdLine,
					Description: context,
					Executed:    false,
					Status:      "pending",
					Timestamp:   time.Time{},
					ExitCode:    -1,
				}

				if cmd.ValidateCommand() {
					commands = append(commands, cmd)
				}
			}
		}
	}

	return commands
}

func isSafeCommand(cmd string) bool {
	tmpCmd := models.Command{
		Raw:    cmd,
		Status: "validating",
	}

	return tmpCmd.ValidateCommand()
}

func executeCommand(cmd models.Command) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(300 * time.Millisecond)

		if !cmd.ValidateCommand() {
			cmd.Status = "blocked"
			cmd.Result = "Command blocked for security reasons"

			return models.CommandResultMsg{
				Command: cmd.Raw,
				Output:  "Command execution blocked for security reasons.",
				Error:   fmt.Errorf("potentially unsafe command"),
			}
		}

		cmd.Status = "executing"
		cmd.Timestamp = time.Now()

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

		if output == "" && err == nil {
			output = "Command executed successfully (no output)"
		}

		if len(output) > 2000 {
			output = output[:2000] + "\n...\n(Output truncated for display)"
		}

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
