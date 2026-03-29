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
