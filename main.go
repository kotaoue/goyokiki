package main

import (
	"fmt"
	"os"
	"time"
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

	filename, err := WriteMarkdownFile(answers, time.Now())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Results written to %s\n", filename)
}
