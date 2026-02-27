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
	got := GenerateMarkdown(answers)
	want := "# 今日やったこと: コードを書いた\n"
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
	got := GenerateMarkdown(answers)
	if !strings.Contains(got, "# 気分はどうですか: よい\n") {
		t.Errorf("missing title line in output: %q", got)
	}
	if !strings.Contains(got, "- [x] よい\n") {
		t.Errorf("missing selected option in output: %q", got)
	}
	if !strings.Contains(got, "- [ ] ふつう\n") {
		t.Errorf("missing unselected option in output: %q", got)
	}
	if !strings.Contains(got, "- [ ] わるい\n") {
		t.Errorf("missing unselected option in output: %q", got)
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
	got := GenerateMarkdown(answers)
	if !strings.HasPrefix(got, "# 今日やったこと: テストを書いた\n") {
		t.Errorf("unexpected start of output: %q", got)
	}
	if !strings.Contains(got, "# 気分: Good\n") {
		t.Errorf("missing single-choice title: %q", got)
	}
	if !strings.Contains(got, "- [x] Good\n") {
		t.Errorf("missing selected option: %q", got)
	}
	if !strings.Contains(got, "- [ ] Bad\n") {
		t.Errorf("missing unselected option: %q", got)
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
	want := "# 今日やったこと: コードを書いた\n"
	if string(content) != want {
		t.Errorf("got %q, want %q", string(content), want)
	}
}
