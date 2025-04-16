package models

import (
	"testing"
)

// TestCommandValidation tests the command validation logic
func TestCommandValidation(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		expected bool
		reason   string
	}{
		// Safe commands - should pass validation
		{"Simple ls", "ls -la", true, ""},
		{"Echo command", "echo 'Hello World'", true, ""},
		{"Touch file", "touch testfile.txt", true, ""},
		{"Create directory", "mkdir -p test/dir", true, ""},
		{"Complex pipe", "find . -type f -name '*.go' | grep 'test' | wc -l", true, ""},
		{"Git command", "git status", true, ""},
		
		// Dangerous commands - should fail validation
		{"Remove root", "rm -rf /", false, "dangerous pattern"},
		{"Remove with wildcard", "rm -rf /*", false, "dangerous pattern"},
		{"Fork bomb", ":(){ :|:& };:", false, "dangerous pattern"},
		{"Fork bomb variant", ":(){:|:&};:", false, "dangerous pattern"},
		{"Format disk", "mkfs /dev/sda1", false, "dangerous pattern"},
		{"Move root files", "mv /* /tmp/", false, "dangerous pattern"},
		{"Change all permissions", "chmod -R 777 /", false, "dangerous pattern"},
		
		// Interactive commands - should fail validation
		{"Vim editor", "vim file.txt", false, "interactive command"},
		{"Nano editor", "nano ~/.bashrc", false, "interactive command"},
		{"Top command", "top", false, "interactive command"},
		{"Less pager", "less /var/log/syslog", false, "interactive command"},
		{"SSH connection", "ssh user@host", false, "interactive command"},
		
		// Command injection - should fail validation
		{"Backtick injection", "echo `rm -rf ~`", false, "command injection"},
		{"Dollar injection", "echo $(rm important-file)", false, "command injection"},
		{"Semicolon injection", "ls; rm -rf ~", false, "command injection"},
		{"AND injection", "ls && rm -rf ~", false, "command injection"},
		
		// Environment variable safety - should fail validation
		{"Set PATH", "export PATH=/malicious:$PATH", false, "sensitive env var"},
		{"Set AWS key", "export AWS_SECRET_ACCESS_KEY=1234", false, "sensitive env var"},
		{"Set password", "export PASSWORD=secret", false, "sensitive env var"},
		
		// Access to sensitive files - should fail validation
		{"Read shadow", "cat /etc/shadow", false, "sensitive file"},
		{"Read passwd", "cat /etc/passwd", false, "sensitive file"},
		{"Read SSH keys", "cat ~/.ssh/id_rsa", false, "sensitive file"},
		
		// Edge cases
		{"Empty command", "", false, "empty command"},
		{"Very long command", "echo " + string(make([]byte, 1000)), false, "too long"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Command{
				Raw: tt.cmd,
			}
			result := cmd.ValidateCommand()
			if result != tt.expected {
				t.Errorf("ValidateCommand() = %v, want %v. Command: %s, Reason: %s", 
					result, tt.expected, tt.cmd, tt.reason)
			}
		})
	}
}
// TestCommandSequenceValidation tests the validation of command sequences
func TestCommandSequenceValidation(t *testing.T) {
	tests := []struct {
		name     string
		cmds     []string
		expected bool
		reason   string
	}{
		{
			"Safe sequence",
			[]string{"mkdir test", "cd test", "touch file.txt", "echo 'hello' > file.txt"},
			true,
			"Basic file operations should be valid",
		},
		{
			"Dangerous file creation and execution",
			[]string{"echo '#!/bin/bash\nrm -rf ~' > script.sh", "chmod +x script.sh", "./script.sh"},
			false,
			"Creating and executing scripts should be blocked",
		},
		{
			"File creation and source",
			[]string{"echo 'export PATH=/bad:$PATH' > setup.sh", "source setup.sh"},
			false,
			"Creating and sourcing env files should be blocked",
		},
		{
			"Multiple safe commands",
			[]string{"ls -la", "pwd", "echo $HOME", "date"},
			true,
			"Simple information commands should be valid",
		},
		{
			"Database manipulation",
			[]string{"mysql -u root -e 'DROP DATABASE production;'"},
			false,
			"Database destruction commands should be blocked",
		},
		{
			"Long dependency chain",
			[]string{
				"mkdir project",
				"cd project",
				"npm init -y",
				"npm install express",
				"echo 'console.log(\"Hello\")' > index.js",
				"node index.js",
			},
			true,
			"Complex but safe project setup should be valid",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var commands []Command
			for _, cmdStr := range tt.cmds {
				commands = append(commands, Command{Raw: cmdStr})
			}
			
			result := ValidateCommandSequence(commands)
			if result != tt.expected {
				t.Errorf("ValidateCommandSequence() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsSystemCommand tests detection of system-level commands
func TestIsSystemCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		expected bool
	}{
		{"sudo command", "sudo apt update", true},
		{"service command", "service nginx restart", true},
		{"systemctl command", "systemctl restart apache2", true},
		{"mount command", "mount /dev/sda1 /mnt", true},
		{"ifconfig command", "ifconfig eth0 down", true},
		
		{"regular command", "ls -la", false},
		{"user file command", "cat file.txt", false},
		{"echo command", "echo hello", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Command{
				Raw: tt.cmd,
			}
			result := cmd.IsSystemCommand()
			if result != tt.expected {
				t.Errorf("IsSystemCommand() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCommandFormatters tests message formatting functions
func TestCommandFormatters(t *testing.T) {
	cmd := Command{
		Raw: "echo 'Hello World'",
		Description: "Print a greeting message",
		Status: "success",
		ExitCode: 0,
	}
	
	// Test ShortDescription
	longDesc := Command{Raw: "test", Description: "This is a very long description that should be truncated when using the ShortDescription method"}
	shortDesc := longDesc.ShortDescription()
	if len(shortDesc) > 50 {
		t.Errorf("ShortDescription() returned string longer than 50 chars: %d", len(shortDesc))
	}
	
	// Test FormattedStatus
	statusCases := []struct {
		status   string
		exitCode int
		expected string
	}{
		{"pending", 0, "Pending"},
		{"executing", 0, "Executing..."},
		{"success", 0, "Succeeded"},
		{"success", 1, "Completed with exit code 1"},
		{"failed", 127, "Failed (exit code: 127)"},
		{"blocked", 0, "Blocked for security reasons"},
		{"unknown", 0, "Unknown status"},
	}
	
	for _, tt := range statusCases {
		cmd.Status = tt.status
		cmd.ExitCode = tt.exitCode
		result := cmd.FormattedStatus()
		if result != tt.expected {
			t.Errorf("FormattedStatus() with status=%s, exitCode=%d = %s, want %s", 
				tt.status, tt.exitCode, result, tt.expected)
		}
	}
}

