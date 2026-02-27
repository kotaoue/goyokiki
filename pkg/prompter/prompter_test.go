package prompter

import (
	"strings"
	"testing"

	"github.com/kotaoue/goyokiki/pkg/questions"
)

func TestPrompter_FreeInput(t *testing.T) {
	input := "doing some work\n"
	q := questions.Question{Title: "What did you do?", Type: questions.FreeInput}
	var sb strings.Builder
	p := NewPrompter(strings.NewReader(input), &sb)
	answers, err := p.Run([]questions.Question{q})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(answers) != 1 {
		t.Fatalf("expected 1 answer, got %d", len(answers))
	}
	if answers[0].Value != "doing some work" {
		t.Errorf("expected %q, got %q", "doing some work", answers[0].Value)
	}
}

func TestPrompter_SingleChoice(t *testing.T) {
	input := "2\n"
	q := questions.Question{
		Title:   "How do you feel?",
		Type:    questions.SingleChoice,
		Options: []string{"Good", "Bad", "Okay"},
	}
	var sb strings.Builder
	p := NewPrompter(strings.NewReader(input), &sb)
	answers, err := p.Run([]questions.Question{q})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if answers[0].Value != "Bad" {
		t.Errorf("expected %q, got %q", "Bad", answers[0].Value)
	}
}

func TestPrompter_SingleChoice_InvalidThenValid(t *testing.T) {
	input := "0\n5\nabc\n1\n"
	q := questions.Question{
		Title:   "Pick one",
		Type:    questions.SingleChoice,
		Options: []string{"A", "B"},
	}
	var sb strings.Builder
	p := NewPrompter(strings.NewReader(input), &sb)
	answers, err := p.Run([]questions.Question{q})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if answers[0].Value != "A" {
		t.Errorf("expected %q, got %q", "A", answers[0].Value)
	}
}
