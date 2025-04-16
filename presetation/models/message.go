package models

import (
	"fmt"
	"strings"
	"time"
)

// Message represents a chat message from either the user or AI
type Message struct {
	Sender   string
	Content  string
	IsAI     bool
	Commands []Command // Commands extracted from AI messages
}

// Command represents a shell command that can be executed
type Command struct {
	Raw         string    // The raw command string
	Description string    // Description or context for the command
	Executed    bool      // Whether the command has been executed
	Timestamp   time.Time // When the command was executed
	Status      string    // Current status: "pending", "executing", "success", "failed", "blocked"
	Result      string    // Output from command execution
	ExitCode    int       // Exit code from command execution
}

// ValidateCommand performs comprehensive validation on a command
func (c Command) ValidateCommand() bool {
	// Check if command is empty
	if c.Raw == "" {
		return false
	}

	// Check if command is too long (reasonable limit)
	if len(c.Raw) > 500 {
		return false
	}

	// 1. Check for dangerous patterns - system damaging commands
	dangerousPatterns := []string{
		"rm -rf /",
		"rm -rf /*",
		"--no-preserve-root",
		":(){:",
		":() {",
		":(){ :|:",
		"f(){ f|f",
		"perl -e 'fork while 1'",
		"./ --no-preserve-root",
		"dd if=/dev/",
		"> /dev/sda",
		"mkfs",
		"mkfs.",
		"mv /* ",
		"chmod -R 777 /",
		":(){ :|:& };:",
		"> /dev/hd",
		"shutdown",
		"halt",
		"poweroff",
		"reboot",
	}

	// Check command against dangerous patterns
	for _, pattern := range dangerousPatterns {
		if strings.Contains(c.Raw, pattern) {
			return false
		}
	}

	// 2. Check for interactive commands that require user input
	interactiveCommands := []string{
		"vim",
		"vi",
		"nano",
		"emacs",
		"pico",
		"top",
		"htop",
		"less",
		"more",
		"man",
		"telnet",
		"ftp",
		"mysql -u",
		"psql -U",
		"ssh ",
		"python -i",
		"ipython",
		"redis-cli",
		"mongo",
		"watch",
	}

	// Check for exact command match or command with arguments
	for _, cmd := range interactiveCommands {
		if c.Raw == cmd || strings.HasPrefix(c.Raw, cmd+" ") {
			return false
		}
	}

	// 3. Check for command injection attempts
	injectionPatterns := []string{
		"$(",
		"`",
		"& ",
		"&&",
		"||",
		"|",
		"; ",
		"< /etc/passwd",
		"> /etc/",
		"eval",
	}

	// 4. Check for suspicious command prefixes that might indicate malicious intent
	suspiciousPrefixes := []string{
		"sudo rm -rf",
		"curl | bash",
		"wget | bash",
		"curl | sh",
		"wget | sh",
		"curl -s | bash",
		"wget -O- | bash",
	}

	for _, prefix := range suspiciousPrefixes {
		if strings.HasPrefix(c.Raw, prefix) {
			return false
		}
	}

	// 5. Special handling for pipe sequences
	// Allow certain safe pipe patterns but validate each component
	if strings.Contains(c.Raw, "|") {
		parts := strings.Split(c.Raw, "|")

		// Check each part of the pipe for dangerous commands
		for _, part := range parts {
			trimmedPart := strings.TrimSpace(part)

			// Skip empty parts
			if trimmedPart == "" {
				continue
			}

			// Block shell commands in pipeline
			if strings.Contains(trimmedPart, "sh ") ||
				strings.Contains(trimmedPart, "bash ") ||
				strings.Contains(trimmedPart, "zsh ") {
				return false
			}

			// Check for command injection in each part
			for _, pattern := range injectionPatterns {
				if strings.Contains(trimmedPart, pattern) {
					return false
				}
			}
		}
	}

	// 6. Environment variable safety checks
	if strings.Contains(c.Raw, "export ") || strings.Contains(c.Raw, "env ") {
		// Block setting PATH environment variable
		if strings.Contains(c.Raw, "PATH=") {
			return false
		}

		// Block setting sensitive environment variables
		sensitiveEnvVars := []string{
			"AWS_", "SECRET_", "PASSWORD", "TOKEN", "KEY", "CREDENTIAL",
		}

		for _, envVar := range sensitiveEnvVars {
			if strings.Contains(strings.ToUpper(c.Raw), envVar) {
				return false
			}
		}
	}

	// 7. Check for commands that access sensitive files
	sensitiveFiles := []string{
		"/etc/shadow",
		"/etc/passwd",
		"/etc/sudoers",
		"/etc/ssh",
		"id_rsa",
		".ssh/",
		".aws/",
		"credentials",
		".env",
	}

	for _, file := range sensitiveFiles {
		if strings.Contains(c.Raw, file) {
			return false
		}
	}

	// If we got here, the command is considered safe
	return true
}

// String returns a string representation of the command
func (c Command) String() string {
	if c.Description != "" {
		return c.Description + ": " + c.Raw
	}
	return c.Raw
}

// IsBlocked checks if the command was blocked for security reasons
func (c Command) IsBlocked() bool {
	return c.Status == "blocked"
}

// IsSuccessful checks if the command executed successfully
func (c Command) IsSuccessful() bool {
	return c.Status == "success" && c.ExitCode == 0
}

// IsPending checks if the command hasn't been executed yet
func (c Command) IsPending() bool {
	return c.Status == "pending" || !c.Executed
}

// IsExecuting checks if the command is currently executing
func (c Command) IsExecuting() bool {
	return c.Status == "executing"
}

// IsFailed checks if the command execution failed
func (c Command) IsFailed() bool {
	return c.Status == "failed" || (c.Executed && c.ExitCode != 0)
}

// Duration calculates the execution time if the command has been executed
// Returns 0 if the command hasn't been executed yet
func (c Command) Duration() time.Duration {
	if !c.Executed || c.Timestamp.IsZero() {
		return 0
	}

	// In a real implementation, we would store both start and end times
	// For now, just return a placeholder duration
	return 1 * time.Second
}

// FormattedStatus returns a user-friendly status string
func (c Command) FormattedStatus() string {
	switch c.Status {
	case "pending":
		return "Pending"
	case "executing":
		return "Executing..."
	case "success":
		if c.ExitCode == 0 {
			return "Succeeded"
		}
		return fmt.Sprintf("Completed with exit code %d", c.ExitCode)
	case "failed":
		return fmt.Sprintf("Failed (exit code: %d)", c.ExitCode)
	case "blocked":
		return "Blocked for security reasons"
	default:
		return "Unknown status"
	}
}

// ShortDescription returns a shortened description for display
func (c Command) ShortDescription() string {
	desc := c.Description
	if len(desc) > 50 {
		desc = desc[:47] + "..."
	}
	return desc
}

// ParseCommandContext extracts context for a command from surrounding text
func ParseCommandContext(content, commandBlock string) string {
	// Find the position of the command in the content
	cmdIndex := strings.Index(content, commandBlock)
	if cmdIndex == -1 {
		return ""
	}

	// Extract text before the command (up to 150 characters)
	startIndex := max(0, cmdIndex-150)
	beforeText := content[startIndex:cmdIndex]

	// Find the last paragraph break or header marker
	breakPoints := []string{"\n\n", "# ", "## ", "### "}
	lastBreakIndex := -1

	for _, breakPoint := range breakPoints {
		if idx := strings.LastIndex(beforeText, breakPoint); idx > lastBreakIndex {
			lastBreakIndex = idx
		}
	}

	if lastBreakIndex != -1 {
		beforeText = beforeText[lastBreakIndex:]
	}

	// Also look for numbered or bulleted list items
	listPatterns := []string{"\n1. ", "\n- ", "\n* ", "\n• "}
	lastListIndex := -1

	for _, pattern := range listPatterns {
		if idx := strings.LastIndex(beforeText, pattern); idx > lastListIndex {
			lastListIndex = idx + 1 // Include the newline
		}
	}

	if lastListIndex > lastBreakIndex {
		beforeText = beforeText[lastListIndex:]
	}

	// Clean up and format the context
	return strings.TrimSpace(beforeText)
}

// ValidateCommandSequence validates a sequence of commands
func ValidateCommandSequence(commands []Command) bool {
	// If no commands, nothing to validate
	if len(commands) == 0 {
		return true
	}

	// Check if individual commands are valid
	for _, cmd := range commands {
		if !cmd.ValidateCommand() {
			return false
		}
	}

	// Check for potentially dangerous sequences
	// For example: creating a file and then executing it
	for i := 0; i < len(commands)-1; i++ {
		// Check for file creation followed by execution
		if strings.Contains(commands[i].Raw, "touch ") ||
			strings.Contains(commands[i].Raw, "echo ") ||
			strings.Contains(commands[i].Raw, "cat >") {

			// Check if next command might execute the file
			if strings.Contains(commands[i+1].Raw, "chmod +x") ||
				strings.Contains(commands[i+1].Raw, "./") ||
				strings.Contains(commands[i+1].Raw, "bash ") ||
				strings.Contains(commands[i+1].Raw, "sh ") {
				return false
			}
		}
	}

	return true
}

// IsSystemCommand determines if a command requires system privileges
func (c Command) IsSystemCommand() bool {
	systemPrefixes := []string{
		"sudo ",
		"doas ",
		"su -c",
		"pkexec ",
		"systemctl ",
		"service ",
	}

	systemCommands := []string{
		"mount",
		"umount",
		"fdisk",
		"fsck",
		"mkfs",
		"chown",
		"passwd",
		"ifconfig",
		"ip addr",
		"ip link",
		"visudo",
	}

	// Check for system prefixes
	for _, prefix := range systemPrefixes {
		if strings.HasPrefix(c.Raw, prefix) {
			return true
		}
	}

	// Check for direct system commands
	for _, cmd := range systemCommands {
		if c.Raw == cmd || strings.HasPrefix(c.Raw, cmd+" ") {
			return true
		}
	}

	return false
}

// FormatErrorMessage formats an error message consistently
func FormatErrorMessage(commandRaw, errorText, output string) string {
	var result strings.Builder

	result.WriteString("## Error executing command\n\n")
	result.WriteString("```bash\n")
	result.WriteString(commandRaw)
	result.WriteString("\n```\n\n")

	result.WriteString("### Error details\n\n")
	result.WriteString(errorText)
	result.WriteString("\n\n")

	if output != "" {
		result.WriteString("### Output\n\n")
		result.WriteString("```\n")
		result.WriteString(output)
		result.WriteString("\n```")
	}

	return result.String()
}

// FormatSuccessMessage formats a success message consistently
func FormatSuccessMessage(commandRaw, output string) string {
	var result strings.Builder

	result.WriteString("## Command executed successfully\n\n")
	result.WriteString("```bash\n")
	result.WriteString(commandRaw)
	result.WriteString("\n```\n\n")

	if output != "" {
		result.WriteString("### Output\n\n")
		result.WriteString("```\n")
		result.WriteString(output)
		result.WriteString("\n```")
	} else {
		result.WriteString("Command completed with no output.")
	}

	return result.String()
}

// FormatCommandPreview generates a preview of the command for confirmation dialog
func FormatCommandPreview(command Command) string {
	var result strings.Builder

	// Format the command itself
	result.WriteString("```bash\n")
	result.WriteString(command.Raw)
	result.WriteString("\n```")

	// Add warning for system commands
	if command.IsSystemCommand() {
		result.WriteString("\n\n**Warning**: This is a system-level command that may require elevated privileges.")
	}

	// Add context if available
	if command.Description != "" {
		result.WriteString("\n\n**Context**: ")
		result.WriteString(command.Description)
	}

	return result.String()
}

// Helper function for ParseCommandContext
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
