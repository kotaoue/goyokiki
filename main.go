package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultConfigPath    = "config.yaml"
	defaultQuestionsPath = "questions.yaml"
)

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

	questionPath := cfg.QuestionFilePath
	if questionPath == "" {
		questionPath = defaultQuestionsPath
	}

	qcfg, err := LoadQuestions(questionPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	prompter := NewPrompter(os.Stdin, os.Stderr)
	answers, err := prompter.Run(qcfg.Questions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	md := GenerateMarkdown(answers)

	if cfg.OutputPath != "" {
		const outputFileLayout = "20060102_150405" // YYYYMMDD_HHMMSS
		filename := time.Now().Format(outputFileLayout) + ".md"
		outPath := filepath.Join(cfg.OutputPath, filename)
		if err := os.WriteFile(outPath, []byte(md), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to write output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Output written to %s\n", outPath)
		return
	}

	fmt.Print(md)
}
