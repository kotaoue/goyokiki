package output

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kotaoue/goyokiki/pkg/prompter"
	"github.com/kotaoue/goyokiki/pkg/questions"
)

func TestGenerateMarkdown_FreeInput(t *testing.T) {
	answers := []prompter.Answer{
		{
			Question: questions.Question{Title: "今日やったこと", Type: questions.FreeInput},
			Value:    "コードを書いた",
		},
	}
	got := GenerateMarkdown("テストタイトル", answers)
	want := "# テストタイトル\n\n- 今日やったこと: コードを書いた\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGenerateMarkdown_SingleChoice(t *testing.T) {
	answers := []prompter.Answer{
		{
			Question: questions.Question{
				Title:   "気分はどうですか",
				Type:    questions.SingleChoice,
				Options: []string{"よい", "ふつう", "わるい"},
			},
			Value: "よい",
		},
	}
	got := GenerateMarkdown("テストタイトル", answers)
	if !strings.Contains(got, "- 気分はどうですか: よい\n") {
		t.Errorf("missing answer line in output: %q", got)
	}
	if strings.Contains(got, "- [x]") || strings.Contains(got, "- [ ]") {
		t.Errorf("output should not contain checkboxes: %q", got)
	}
}

func TestGenerateMarkdown_Mixed(t *testing.T) {
	answers := []prompter.Answer{
		{
			Question: questions.Question{Title: "今日やったこと", Type: questions.FreeInput},
			Value:    "テストを書いた",
		},
		{
			Question: questions.Question{
				Title:   "気分",
				Type:    questions.SingleChoice,
				Options: []string{"Good", "Bad"},
			},
			Value: "Good",
		},
	}
	got := GenerateMarkdown("テストタイトル", answers)
	if !strings.HasPrefix(got, "# テストタイトル\n\n") {
		t.Errorf("unexpected start of output: %q", got)
	}
	if !strings.Contains(got, "- 今日やったこと: テストを書いた\n") {
		t.Errorf("missing free input line: %q", got)
	}
	if !strings.Contains(got, "- 気分: Good\n") {
		t.Errorf("missing single-choice line: %q", got)
	}
	if strings.Contains(got, "- [x]") || strings.Contains(got, "- [ ]") {
		t.Errorf("output should not contain checkboxes: %q", got)
	}
}

func TestResolveFilename(t *testing.T) {
	answers := []prompter.Answer{
		{
			Question: questions.Question{Title: "お店", Type: questions.FreeInput},
			Value:    "TAKAO COFFEE",
		},
		{
			Question: questions.Question{Title: "メニュー名", Type: questions.FreeInput},
			Value:    "カフェラテ",
		},
	}

	tests := []struct {
		name     string
		template string
		want     string
	}{
		{
			name:     "single marker",
			template: "{ANSWER_1}",
			want:     "TAKAO COFFEE",
		},
		{
			name:     "two markers",
			template: "{ANSWER_1}:{ANSWER_2}",
			want:     "TAKAO COFFEE:カフェラテ",
		},
		{
			name:     "out-of-range marker is left unchanged",
			template: "{ANSWER_1}:{ANSWER_99}",
			want:     "TAKAO COFFEE:{ANSWER_99}",
		},
		{
			name:     "no markers",
			template: "fixed-name",
			want:     "fixed-name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveFilename(tt.template, answers)
			if got != tt.want {
				t.Errorf("ResolveFilename(%q) = %q, want %q", tt.template, got, tt.want)
			}
		})
	}
}

func TestGenerateMarkdown_Empty(t *testing.T) {
	got := GenerateMarkdown("空タイトル", nil)
	want := "# 空タイトル\n\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteMarkdownFile(t *testing.T) {
	answers := []prompter.Answer{
		{
			Question: questions.Question{Title: "今日やったこと", Type: questions.FreeInput},
			Value:    "コードを書いた",
		},
	}
	now := time.Date(2026, 2, 21, 13, 25, 33, 0, time.UTC)
	filename, err := WriteMarkdownFile(answers, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.Remove(filename)

	if filename != "results-20260221132533.md" {
		t.Errorf("unexpected filename: %q", filename)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	want := "# results-20260221132533\n\n- 今日やったこと: コードを書いた\n"
	if string(content) != want {
		t.Errorf("got %q, want %q", string(content), want)
	}
}

func TestWriteMarkdownFile_Error(t *testing.T) {
	// Writing to a path whose parent directory does not exist must return an error.
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	// Switch to a read-only temp dir so the write is denied.
	roDir := t.TempDir()
	if err := os.Chmod(roDir, 0555); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	if err := os.Chdir(roDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	answers := []prompter.Answer{
		{
			Question: questions.Question{Title: "テスト", Type: questions.FreeInput},
			Value:    "値",
		},
	}
	_, writeErr := WriteMarkdownFile(answers, time.Now())
	if writeErr == nil {
		t.Error("expected an error when writing to a read-only directory, got nil")
	}
}
