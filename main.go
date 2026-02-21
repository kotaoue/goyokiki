package main

import (
	"fmt"
	"os"
)

const defaultConfigPath = "config.yaml"

func main() {
	configPath := defaultConfigPath
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	prompter := NewPrompter(os.Stdin, os.Stderr)
	answers, err := prompter.Run(cfg.Questions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(GenerateMarkdown(answers))
}
