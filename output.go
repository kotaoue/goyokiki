package main

import (
	"fmt"
	"strings"
)

// GenerateMarkdown converts a slice of answers into a Markdown string.
func GenerateMarkdown(answers []Answer) string {
	var sb strings.Builder
	for i, a := range answers {
		if i > 0 {
			sb.WriteString("\n")
		}
		switch a.Question.Type {
		case FreeInput:
			fmt.Fprintf(&sb, "# %s: %s\n", a.Question.Title, a.Value)
		case SingleChoice:
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
