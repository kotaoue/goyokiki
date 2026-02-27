package output

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kotaoue/goyokiki/pkg/prompter"
	"github.com/kotaoue/goyokiki/pkg/questions"
)

// GenerateMarkdown converts a slice of answers into a Markdown string.
func GenerateMarkdown(answers []prompter.Answer) string {
	var sb strings.Builder
	for i, a := range answers {
		if i > 0 {
			sb.WriteString("\n")
		}
		switch a.Question.Type {
		case questions.FreeInput:
			fmt.Fprintf(&sb, "# %s: %s\n", a.Question.Title, a.Value)
		case questions.SingleChoice:
			fmt.Fprintf(&sb, "# %s: %s\n", a.Question.Title, a.Value)
			for _, opt := range a.Question.Options {
				if opt == a.Value {
					fmt.Fprintf(&sb, "- [x] %s\n", opt)
				} else {
					fmt.Fprintf(&sb, "- [ ] %s\n", opt)
				}
			}
		}
	}
	return sb.String()
}

// WriteMarkdownFile writes the Markdown output to a file named results-yyyymmddhhiiss.md.
// It returns the filename that was written.
func WriteMarkdownFile(answers []prompter.Answer, now time.Time) (string, error) {
	filename := fmt.Sprintf("results-%s.md", now.Format("20060102150405"))
	if err := os.WriteFile(filename, []byte(GenerateMarkdown(answers)), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	return filename, nil
}
