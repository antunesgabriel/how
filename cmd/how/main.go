package main

import (
	"fmt"
	"os"

	"github.com/antunesgabriel/how-ai/presetation"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	chat, err := presetation.NewChat()
	if err != nil {
		panic(err)
	}

	p := tea.NewProgram(chat)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
