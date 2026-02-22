package goyokiki

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
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

func TestLoadQuestions(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		path      string // if set, used directly instead of a temp file
		wantErr   bool
		wantCount int
	}{
		{
			name: "valid",
			content: `
questions:
  - title: "What did you do?"
    type: free
  - title: "How do you feel?"
    type: single
    options:
      - Good
      - Bad
`,
			wantCount: 2,
		},
		{
			name:    "file not found",
			path:    "/nonexistent/path/questions.yaml",
			wantErr: true,
		},
		{
			name:    "missing title",
			content: "questions:\n  - type: free\n",
			wantErr: true,
		},
		{
			name:    "unknown type",
			content: "questions:\n  - title: \"Q\"\n    type: unknown\n",
			wantErr: true,
		},
		{
			name:    "single choice no options",
			content: "questions:\n  - title: \"Q\"\n    type: single\n",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.path
			if p == "" {
				p = writeTempFile(t, tc.content)
			}
			qcfg, err := LoadQuestions(p)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantCount > 0 && len(qcfg.Questions) != tc.wantCount {
				t.Errorf("expected %d questions, got %d", tc.wantCount, len(qcfg.Questions))
			}
		})
	}
}

func TestLoadConfig_WithQuestionFilePath(t *testing.T) {
	questions := `
questions:
  - title: "External question"
    type: free
`
	questionPath := writeTempFile(t, questions)

	content := fmt.Sprintf(`question_file: %s`, questionPath)
	cfgPath := writeTempFile(t, content)

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.QuestionFilePath != questionPath {
		t.Errorf("unexpected question_file: %q", cfg.QuestionFilePath)
	}

	qcfg, err := LoadQuestions(cfg.QuestionFilePath)
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
