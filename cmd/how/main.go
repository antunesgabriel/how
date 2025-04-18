package main

import (
	"fmt"
	"os"

	"github.com/antunesgabriel/how-ai/presetation"
)

func main() {
	if err := presetation.StartApp(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
