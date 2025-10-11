package bootstrap

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type PromptDecision string

const (
	PromptKeep      PromptDecision = "keep"
	PromptOverwrite PromptDecision = "overwrite"
	PromptCancel    PromptDecision = "cancel"
)

type PromptOutcome struct {
	Decision PromptDecision
}

func PromptForOverwrite(ioCfg IO, path string) (PromptOutcome, error) {
	ioCfg = normalizeIO(ioCfg)
	writer := ioCfg.Stderr
	reader := bufio.NewReader(ioCfg.Stdin)

	fmt.Fprintf(writer, "sysgreet bootstrap: configuration already exists at %s\n", path)
	fmt.Fprintln(writer, "Choose an option:")
	fmt.Fprintln(writer, "  [K]eep existing config")
	fmt.Fprintln(writer, "  [O]verwrite with defaults (backup will be created)")
	fmt.Fprintln(writer, "  [C]ancel and exit")

	for {
		fmt.Fprint(writer, "Selection [K/O/C]: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return PromptOutcome{Decision: PromptCancel}, nil
			}
			return PromptOutcome{}, fmt.Errorf("prompt read: %w", err)
		}
		choice := strings.ToLower(strings.TrimSpace(line))
		if choice == "k" || choice == "keep" {
			return PromptOutcome{Decision: PromptKeep}, nil
		}
		if choice == "o" || choice == "overwrite" {
			return PromptOutcome{Decision: PromptOverwrite}, nil
		}
		if choice == "c" || choice == "cancel" {
			return PromptOutcome{Decision: PromptCancel}, nil
		}
		fmt.Fprintln(writer, "Invalid selection. Please choose K, O, or C.")
	}
}
