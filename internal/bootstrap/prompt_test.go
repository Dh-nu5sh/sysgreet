package bootstrap

import (
	"bytes"
	"strings"
	"testing"
)

func TestPromptForOverwriteKeepsConfig(t *testing.T) {
	input := bytes.NewBufferString("k\n")
	stderr := &bytes.Buffer{}

	outcome, err := PromptForOverwrite(IO{Stdin: input, Stderr: stderr}, "/tmp/config.yaml")
	if err != nil {
		t.Fatalf("PromptForOverwrite error: %v", err)
	}
	if outcome.Decision != PromptKeep {
		t.Fatalf("expected keep decision, got %s", outcome.Decision)
	}
	if !strings.Contains(stderr.String(), "[K]eep") {
		t.Fatalf("expected prompt output to show options, got %q", stderr.String())
	}
}

func TestPromptForOverwriteHandlesInvalidInput(t *testing.T) {
	input := bytes.NewBufferString("x\no\n")
	stderr := &bytes.Buffer{}

	outcome, err := PromptForOverwrite(IO{Stdin: input, Stderr: stderr}, "config.yaml")
	if err != nil {
		t.Fatalf("PromptForOverwrite error: %v", err)
	}
	if outcome.Decision != PromptOverwrite {
		t.Fatalf("expected overwrite decision, got %s", outcome.Decision)
	}
	if count := strings.Count(stderr.String(), "Invalid selection"); count == 0 {
		t.Fatalf("expected invalid selection message in prompt output: %q", stderr.String())
	}
}

func TestPromptForOverwriteCancel(t *testing.T) {
	input := bytes.NewBufferString("C\n")
	stderr := &bytes.Buffer{}

	outcome, err := PromptForOverwrite(IO{Stdin: input, Stderr: stderr}, "config.yaml")
	if err != nil {
		t.Fatalf("PromptForOverwrite error: %v", err)
	}
	if outcome.Decision != PromptCancel {
		t.Fatalf("expected cancel decision, got %s", outcome.Decision)
	}
}
