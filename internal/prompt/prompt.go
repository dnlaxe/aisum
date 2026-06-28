package prompt

import (
	"fmt"
	"strings"
)

func TerminalSummary(input string) string {
	input = strings.TrimSpace(input)

	return fmt.Sprintf(`You are a command-line assistant.

Summarize the following terminal output for a developer.

Terminal output:

%s`, input)
}
