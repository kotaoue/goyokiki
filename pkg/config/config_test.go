package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/kotaoue/goyokiki/pkg/questions"
)

func TestLoadConfig_Valid(t *testing.T) {
	content := `
question_file: /path/to/questions.yaml
output_path: /path/to/output/
`
	path := writeTempFile(t, content)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.QuestionFilePath != "/path/to/questions.yaml" {
		t.Errorf("unexpected question_file: %q", cfg.QuestionFilePath)
	}
	if cfg.OutputPath != "/path/to/output/" {
		t.Errorf("unexpected output_path: %q", cfg.OutputPath)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_WithQuestionFilePath(t *testing.T) {
	qs := `
questions:
  - title: "External question"
    type: free
`
	questionPath := writeTempFile(t, qs)

	content := fmt.Sprintf(`question_file: %s`, questionPath)
	cfgPath := writeTempFile(t, content)

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.QuestionFilePath != questionPath {
		t.Errorf("unexpected question_file: %q", cfg.QuestionFilePath)
	}

	qcfg, err := questions.LoadQuestions(cfg.QuestionFilePath)
	if err != nil {
		t.Fatalf("unexpected error loading questions: %v", err)
	}
	if len(qcfg.Questions) != 1 {
		t.Fatalf("expected 1 question, got %d", len(qcfg.Questions))
	}
	if qcfg.Questions[0].Title != "External question" {
		t.Errorf("unexpected question title: %q", qcfg.Questions[0].Title)
	}
}

func TestLoadConfig_WithOutputPath(t *testing.T) {
	dir := t.TempDir()
	content := fmt.Sprintf(`output_path: %s`, dir)
	cfgPath := writeTempFile(t, content)
	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputPath != dir {
		t.Errorf("expected OutputPath %q, got %q", dir, cfg.OutputPath)
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
