package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Valid(t *testing.T) {
	content := `
questions:
  - title: "What did you do?"
    type: free
  - title: "How do you feel?"
    type: single
    options:
      - Good
      - Bad
`
	path := writeTempFile(t, content)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Questions) != 2 {
		t.Fatalf("expected 2 questions, got %d", len(cfg.Questions))
	}
	if cfg.Questions[0].Type != FreeInput {
		t.Errorf("expected FreeInput, got %q", cfg.Questions[0].Type)
	}
	if cfg.Questions[1].Type != SingleChoice {
		t.Errorf("expected SingleChoice, got %q", cfg.Questions[1].Type)
	}
	if len(cfg.Questions[1].Options) != 2 {
		t.Errorf("expected 2 options, got %d", len(cfg.Questions[1].Options))
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_MissingTitle(t *testing.T) {
	content := `
questions:
  - type: free
`
	path := writeTempFile(t, content)
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected validation error for missing title")
	}
}

func TestLoadConfig_UnknownType(t *testing.T) {
	content := `
questions:
  - title: "Question"
    type: unknown
`
	path := writeTempFile(t, content)
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected validation error for unknown type")
	}
}

func TestLoadConfig_SingleChoiceNoOptions(t *testing.T) {
	content := `
questions:
  - title: "Question"
    type: single
`
	path := writeTempFile(t, content)
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected validation error for single-choice with no options")
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}
