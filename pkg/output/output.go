package output

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kotaoue/goyokiki/pkg/prompter"
)

// GenerateMarkdown converts a slice of answers into a Markdown string with
// a single top-level heading followed by a bullet-point list of all answers.
func GenerateMarkdown(title string, answers []prompter.Answer) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "# %s\n\n", title)
	for _, a := range answers {
		fmt.Fprintf(&sb, "- %s: %s\n", a.Question.Title, a.Value)
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
	title := fmt.Sprintf("results-%s", now.Format("20060102150405"))
	filename := title + ".md"
	if err := os.WriteFile(filename, []byte(GenerateMarkdown(title, answers)), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	return filename, nil
}
