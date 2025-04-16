package mock

var AIResponses = []string{
	`# Hello there! 👋

I'm your AI assistant. How can I help you today?`,

	`## Docker Installation on macOS

To install Docker on your MacBook, follow these steps:

1. Visit [Docker's official website](https://www.docker.com/products/docker-desktop/)
2. Download Docker Desktop for Mac
3. Open the .dmg file and drag Docker to your Applications folder
4. Open Docker from your Applications folder
5. Wait for Docker to start (you'll see the whale icon in your menu bar)

**System Requirements:**
- macOS 11 or newer
- At least 4GB of RAM

Would you like to know how to verify your installation?`,

	`# Go Programming Tips

Here are some useful Go packages for CLI applications:

| Package | Purpose | URL |
|---------|---------|-----|
| Cobra | Command-line interface | github.com/spf13/cobra |
| Viper | Configuration | github.com/spf13/viper |
| BubbleTea | Terminal UI | github.com/charmbracelet/bubbletea |

Go's standard library is also quite powerful for CLI tools!

` + "```go" + `
package main

import "fmt"

func main() {
    fmt.Println("Hello, Gopher!")
}
` + "```",
}
