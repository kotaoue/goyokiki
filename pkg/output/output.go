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

// ResolveFilename replaces {ANSWER_N} markers in the template with the
// corresponding answer value (1-indexed). Markers with an out-of-range
// index are left unchanged.
func ResolveFilename(template string, answers []prompter.Answer) string {
	result := template
	for i, a := range answers {
		marker := fmt.Sprintf("{ANSWER_%d}", i+1)
		result = strings.ReplaceAll(result, marker, a.Value)
	}
	return result
}

// WriteMarkdownFile writes the Markdown output to a file named results-yyyymmddhhiiss.md.
func WriteMarkdownFile(answers []prompter.Answer, now time.Time) (string, error) {
	filename := fmt.Sprintf("results-%s.md", now.Format("20060102150405"))
	if err := os.WriteFile(filename, []byte(GenerateMarkdown(answers)), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	return filename, nil
}
