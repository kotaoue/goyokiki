package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kotaoue/goyokiki/pkg/config"
	"github.com/kotaoue/goyokiki/pkg/output"
	"github.com/kotaoue/goyokiki/pkg/prompter"
	"github.com/kotaoue/goyokiki/pkg/questions"
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

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	questionPath := cfg.QuestionFilePath
	if questionPath == "" {
		questionPath = defaultQuestionsPath
	}

	qcfg, err := questions.LoadQuestions(questionPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	p := prompter.NewPrompter(os.Stdin, os.Stderr)
	answers, err := p.Run(qcfg.Questions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var title string
	if cfg.OutputFilename != "" {
		title = output.ResolveFilename(cfg.OutputFilename, answers)
	} else {
		const outputFileLayout = "20060102_150405" // YYYYMMDD_HHMMSS
		title = time.Now().Format(outputFileLayout)
	}

	md := output.GenerateMarkdown(title, answers)

	if cfg.OutputPath != "" {
		outPath := filepath.Join(cfg.OutputPath, title+".md")
		if err := os.WriteFile(outPath, []byte(md), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to write output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Output written to %s\n", outPath)
		return
	}

	fmt.Print(md)
}
