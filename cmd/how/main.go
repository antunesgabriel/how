package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/antunesgabriel/how-ai/config"
	"github.com/antunesgabriel/how-ai/presetation"
)

func main() {
	if len(os.Args) > 1 {
		cmd := os.Args[1]

		if cmd == "init" {
			if err := handleInit(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		if cmd == "example" {
			if err := config.CreateExampleConfig(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Example configuration created at %s\n", config.ConfigFilePath())
			fmt.Println("Please edit this file with your API keys and preferences.")
			return
		}

		query := strings.Join(os.Args[1:], " ")
		if err := startApp(query); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := startApp(""); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func handleInit() error {
	configPath := config.ConfigFilePath()
	_, err := os.Stat(configPath)
	if err == nil {
		fmt.Printf("Configuration file already exists at %s\n", configPath)
		fmt.Println("To create an example configuration with all providers, run: how example")
		return nil
	}

	if err := config.CreateDefaultConfig(); err != nil {
		return err
	}

	fmt.Printf("Default configuration created at %s\n", configPath)
	fmt.Println("Please edit this file with your API key and preferences.")
	return nil
}

func startApp(query string) error {
	cfg, err := config.Load()
	if err != nil {
		if strings.Contains(err.Error(), "config file not found") {
			fmt.Println("Configuration file not found.")
			fmt.Println("Run 'how init' to create a default configuration.")
			return fmt.Errorf("configuration required")
		}
		return err
	}

	if err := presetation.StartApp(cfg, query); err != nil {
		return err
	}

	return nil
}
