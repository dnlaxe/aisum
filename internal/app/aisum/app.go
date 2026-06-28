package aisum

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/dnlaxe/aisum/internal/ollama"
	"github.com/dnlaxe/aisum/internal/prompt"
)

func Run(stdin io.Reader, stdout io.Writer, stderr io.Writer) int {

	input, err := io.ReadAll(stdin)
	if err != nil {
		fmt.Fprintf(stderr, "aisum: read stdin: %v\n", err)
		return 1
	}

	input = bytes.TrimSpace(input)
	if len(input) == 0 {
		fmt.Println(stderr, "aisum: no input provided")
		return 1
	}

	requestPrompt := prompt.TerminalSummary(string(input))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client := ollama.NewClient(
		ollama.DefaultBaseUrl,
		ollama.DefaultModel,
	)

	output, err := client.Generate(ctx, requestPrompt)
	if err != nil {
		fmt.Fprintf(stderr, "aisum: %v\n", err)
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "Make sure Ollama is running and the model is installed:")
		fmt.Fprintln(stderr, "  ollama serve")
		fmt.Fprintf(stderr, "  ollama pull %s\n", ollama.DefaultModel)
		return 1
	}

	if _, err := fmt.Fprintln(stdout, output); err != nil {
		fmt.Fprintf(stderr, "aisum: write stdout: %v\n", err)
		return 1
	}

	return 0
}
