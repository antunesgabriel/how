package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/antunesgabriel/how-ai/presetation"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	rand.Seed(time.Now().UnixNano())

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
