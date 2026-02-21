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

// Config holds the full configuration loaded from a YAML file.
type Config struct {
	Questions []Question `yaml:"questions"`
}

// LoadConfig reads and parses a YAML configuration file at the given path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	for i, q := range c.Questions {
		if q.Title == "" {
			return fmt.Errorf("question %d: title is required", i+1)
		}
		switch q.Type {
		case FreeInput:
			// no extra fields required
		case SingleChoice:
			if len(q.Options) == 0 {
				return fmt.Errorf("question %d (%q): single-choice question requires at least one option", i+1, q.Title)
			}
		default:
			return fmt.Errorf("question %d (%q): unknown type %q (must be %q or %q)", i+1, q.Title, q.Type, FreeInput, SingleChoice)
		}
	}
	return nil
}
