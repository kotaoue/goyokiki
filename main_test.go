package main

import (
	"testing"
	"time"

	"github.com/kotaoue/goyokiki/pkg/prompter"
	"github.com/kotaoue/goyokiki/pkg/questions"
)

func TestResolveTitle_WithTemplate(t *testing.T) {
	answers := []prompter.Answer{
		{Question: questions.Question{Title: "お店", Type: questions.FreeInput}, Value: "TAKAO COFFEE"},
		{Question: questions.Question{Title: "メニュー名", Type: questions.FreeInput}, Value: "カフェラテ"},
	}
	now := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)

	got := resolveTitle("{ANSWER_1}:{ANSWER_2}", answers, now)
	want := "TAKAO COFFEE:カフェラテ"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolveTitle_NoTemplate_UsesTimestamp(t *testing.T) {
	answers := []prompter.Answer{
		{Question: questions.Question{Title: "メモ", Type: questions.FreeInput}, Value: "てすと"},
	}
	now := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)

	got := resolveTitle("", answers, now)
	want := "20260102_150405"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolveTitle_TemplateOutOfRange(t *testing.T) {
	answers := []prompter.Answer{
		{Question: questions.Question{Title: "お店", Type: questions.FreeInput}, Value: "TAKAO COFFEE"},
	}
	now := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)

	// {ANSWER_99} is out of range and should be left unchanged
	got := resolveTitle("{ANSWER_1}-{ANSWER_99}", answers, now)
	want := "TAKAO COFFEE-{ANSWER_99}"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
