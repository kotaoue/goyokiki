package main

import (
	"strings"
	"testing"
)

func TestPrompter_FreeInput(t *testing.T) {
	input := "doing some work\n"
	q := Question{Title: "What did you do?", Type: FreeInput}
	var sb strings.Builder
	p := NewPrompter(strings.NewReader(input), &sb)
	answers, err := p.Run([]Question{q})
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
	q := Question{
		Title:   "How do you feel?",
		Type:    SingleChoice,
		Options: []string{"Good", "Bad", "Okay"},
	}
	var sb strings.Builder
	p := NewPrompter(strings.NewReader(input), &sb)
	answers, err := p.Run([]Question{q})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if answers[0].Value != "Bad" {
		t.Errorf("expected %q, got %q", "Bad", answers[0].Value)
	}
}

func TestPrompter_SingleChoice_InvalidThenValid(t *testing.T) {
	input := "0\n5\nabc\n1\n"
	q := Question{
		Title:   "Pick one",
		Type:    SingleChoice,
		Options: []string{"A", "B"},
	}
	var sb strings.Builder
	p := NewPrompter(strings.NewReader(input), &sb)
	answers, err := p.Run([]Question{q})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if answers[0].Value != "A" {
		t.Errorf("expected %q, got %q", "A", answers[0].Value)
	}
}
