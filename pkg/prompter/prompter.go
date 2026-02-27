package prompter

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/kotaoue/goyokiki/pkg/questions"
)

// Answer holds the user's answer to a question.
type Answer struct {
	Question questions.Question
	Value    string // selected text or free-text input
}

// Prompter handles interactive user prompts.
type Prompter struct {
	in  *bufio.Reader
	out io.Writer
}

// NewPrompter creates a Prompter that reads from r and writes to w.
func NewPrompter(r io.Reader, w io.Writer) *Prompter {
	return &Prompter{in: bufio.NewReader(r), out: w}
}

// Run iterates through all questions, prompts the user, and returns answers.
func (p *Prompter) Run(qs []questions.Question) ([]Answer, error) {
	answers := make([]Answer, 0, len(qs))
	for _, q := range qs {
		var ans Answer
		var err error
		switch q.Type {
		case questions.FreeInput:
			ans, err = p.promptFree(q)
		case questions.SingleChoice:
			ans, err = p.promptSingle(q)
		}
		if err != nil {
			return nil, err
		}
		answers = append(answers, ans)
	}
	return answers, nil
}

func (p *Prompter) promptFree(q questions.Question) (Answer, error) {
	fmt.Fprintf(p.out, "%s: ", q.Title)
	line, err := p.in.ReadString('\n')
	if err != nil && err != io.EOF {
		return Answer{}, fmt.Errorf("failed to read input: %w", err)
	}
	return Answer{Question: q, Value: strings.TrimRight(line, "\r\n")}, nil
}

func (p *Prompter) promptSingle(q questions.Question) (Answer, error) {
	fmt.Fprintf(p.out, "%s\n", q.Title)
	for i, opt := range q.Options {
		fmt.Fprintf(p.out, "  %d) %s\n", i+1, opt)
	}
	for {
		fmt.Fprintf(p.out, "Enter choice (1-%d): ", len(q.Options))
		line, err := p.in.ReadString('\n')
		if err != nil && err != io.EOF {
			return Answer{}, fmt.Errorf("failed to read input: %w", err)
		}
		line = strings.TrimRight(line, "\r\n")
		n, convErr := strconv.Atoi(strings.TrimSpace(line))
		if convErr == nil && n >= 1 && n <= len(q.Options) {
			return Answer{Question: q, Value: q.Options[n-1]}, nil
		}
		fmt.Fprintf(p.out, "Invalid choice. Please enter a number between 1 and %d.\n", len(q.Options))
	}
}
