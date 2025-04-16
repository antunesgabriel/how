package mock

// Test fixture for AI responses with commands
var TestAIResponses = map[string]string{
	"welcome": "# Hello there!\n\nI'm your friendly terminal AI assistant. How can I help you today?",

	"file_commands": "# File System Operations\n\nHere are some common file system operations:\n\n- List files: `ls -la`\n- Create directory: `mkdir dirname`\n- Change directory: `cd dirname`\n- Remove file: `rm filename`\n- Copy file: `cp source dest`",

	"golang_commands": "# Golang Development\n\nTo build a Go application:\n\n```bash\ngo build main.go\n```\n\nTo run tests:\n\n```bash\ngo test ./...\n```",

	"system_info": "# Checking System Information\n\nYou can check your system information using this command:\n\n```bash\nuname -a\n```\n\nAnd list your current directory contents with:\n\n```bash\nls -la\n```",

	"dangerous_command": "# Warning: Dangerous Command\n\nNever run this command as it would delete your entire filesystem:\n\n```bash\nrm -rf /\n```\n\nInstead, use targeted removal for specific files.",

	"multi_command": "# Setting Up a New Project\n\nHere's how to set up a new project:\n\n```bash\nmkdir -p myproject\ncd myproject\ntouch README.md\necho '# My Project' > README.md\ngit init\n```\n\nThis creates a basic project structure.",
}
