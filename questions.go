package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// QuestionType represents the type of a question.
type QuestionType string

const (
	// FreeInput represents a free-text input question.
	FreeInput QuestionType = "free"
	// SingleChoice represents a single-choice selection question.
	SingleChoice QuestionType = "single"
)

// Question holds the definition of a single prompt question.
type Question struct {
	Title   string       `yaml:"title"`
	Type    QuestionType `yaml:"type"`
	Options []string     `yaml:"options,omitempty"`
}

// QuestionsConfig holds the questions loaded from a YAML questions file.
type QuestionsConfig struct {
	Questions []Question `yaml:"questions"`
}

// LoadQuestions reads and parses a YAML questions file at the given path.
func LoadQuestions(path string) (*QuestionsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read questions file: %w", err)
	}

	var qcfg QuestionsConfig
	if err := yaml.Unmarshal(data, &qcfg); err != nil {
		return nil, fmt.Errorf("failed to parse questions file: %w", err)
	}

	if err := qcfg.validate(); err != nil {
		return nil, err
	}

	return &qcfg, nil
}

func (q *QuestionsConfig) validate() error {
	for i, question := range q.Questions {
		if question.Title == "" {
			return fmt.Errorf("question %d: title is required", i+1)
		}
		switch question.Type {
		case FreeInput:
			// no extra fields required
		case SingleChoice:
			if len(question.Options) == 0 {
				return fmt.Errorf("question %d (%q): single-choice question requires at least one option", i+1, question.Title)
			}
		default:
			return fmt.Errorf("question %d (%q): unknown type %q (must be %q or %q)", i+1, question.Title, question.Type, FreeInput, SingleChoice)
		}
	}
	return nil
}
