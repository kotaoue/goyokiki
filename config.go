package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the path settings loaded from a YAML config file.
type Config struct {
	QuestionFilePath string `yaml:"question_file"`
	OutputPath       string `yaml:"output_path,omitempty"`
}

// LoadConfig reads and parses a YAML config file at the given path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}
