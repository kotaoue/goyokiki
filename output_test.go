package main

import (
	"strings"
	"testing"
)

func TestGenerateMarkdown_FreeInput(t *testing.T) {
	answers := []Answer{
		{
			Question: Question{Title: "今日やったこと", Type: FreeInput},
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
	answers := []Answer{
		{
			Question: Question{
				Title:   "気分はどうですか",
				Type:    SingleChoice,
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
	answers := []Answer{
		{
			Question: Question{Title: "今日やったこと", Type: FreeInput},
			Value:    "テストを書いた",
		},
		{
			Question: Question{
				Title:   "気分",
				Type:    SingleChoice,
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
