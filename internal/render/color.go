package render

import (
	"os"
	"strings"
)

var ansiCodes = map[string]string{
	"red":    "\033[31m",
	"yellow": "\033[33m",
	"green":  "\033[32m",
	"cyan":   "\033[36m",
	"reset":  "\033[0m",
}

// Colorizer wraps text in ANSI sequences when enabled.
type Colorizer struct {
	enabled bool
}

// NewColorizer creates a colorizer that may be disabled via NO_COLOR or explicit flag.
func NewColorizer(disable bool) Colorizer {
	if disable {
		return Colorizer{enabled: false}
	}
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return Colorizer{enabled: false}
	}
	return Colorizer{enabled: true}
}

// Wrap applies color when enabled.
func (c Colorizer) Wrap(color, text string) string {
	if !c.enabled {
		return text
	}
	code, ok := ansiCodes[color]
	if !ok || color == "reset" {
		return text
	}
	return code + text + ansiCodes["reset"]
}

// Strip removes ANSI escape sequences from a string.
func Strip(input string) string {
	var b strings.Builder
	i := 0
	runes := []rune(input)
	for i < len(runes) {
		if runes[i] == '' {
			// Skip until letter
			i++
			for i < len(runes) && ((runes[i] >= '0' && runes[i] <= '9') || runes[i] == '[' || runes[i] == ';') {
				i++
			}
			if i < len(runes) {
				i++
			}
			continue
		}
		b.WriteRune(runes[i])
		i++
	}
	return b.String()
}
