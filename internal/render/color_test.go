package render

import (
	"os"
	"strings"
	"testing"
)

func TestNewColorizer(t *testing.T) {
	tests := []struct {
		name        string
		disable     bool
		setNoColor  bool
		wantEnabled bool
	}{
		{
			name:        "enabled by default",
			disable:     false,
			setNoColor:  false,
			wantEnabled: true,
		},
		{
			name:        "explicitly disabled",
			disable:     true,
			setNoColor:  false,
			wantEnabled: false,
		},
		{
			name:        "disabled via NO_COLOR env",
			disable:     false,
			setNoColor:  true,
			wantEnabled: false,
		},
		{
			name:        "both disabled",
			disable:     true,
			setNoColor:  true,
			wantEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup NO_COLOR environment
			if tt.setNoColor {
				os.Setenv("NO_COLOR", "1")
				defer os.Unsetenv("NO_COLOR")
			}

			c := NewColorizer(tt.disable)
			if c.enabled != tt.wantEnabled {
				t.Errorf("NewColorizer(%v).enabled = %v, want %v", tt.disable, c.enabled, tt.wantEnabled)
			}
		})
	}
}

func TestColorizer_Wrap(t *testing.T) {
	tests := []struct {
		name         string
		enabled      bool
		color        string
		text         string
		wantContains []string
		wantEqual    bool
	}{
		{
			name:         "red color enabled",
			enabled:      true,
			color:        "red",
			text:         "ERROR",
			wantContains: []string{"\033[31m", "ERROR", "\033[0m"},
		},
		{
			name:         "yellow color enabled",
			enabled:      true,
			color:        "yellow",
			text:         "WARNING",
			wantContains: []string{"\033[33m", "WARNING", "\033[0m"},
		},
		{
			name:         "green color enabled",
			enabled:      true,
			color:        "green",
			text:         "OK",
			wantContains: []string{"\033[32m", "OK", "\033[0m"},
		},
		{
			name:      "cyan color enabled",
			enabled:   true,
			color:     "cyan",
			text:      "INFO",
			wantContains: []string{"\033[36m", "INFO", "\033[0m"},
		},
		{
			name:      "disabled colorizer",
			enabled:   false,
			color:     "red",
			text:      "ERROR",
			wantEqual: true,
		},
		{
			name:      "reset color is passthrough",
			enabled:   true,
			color:     "reset",
			text:      "TEXT",
			wantEqual: true,
		},
		{
			name:      "unknown color is passthrough",
			enabled:   true,
			color:     "magenta",
			text:      "TEXT",
			wantEqual: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Colorizer{enabled: tt.enabled}
			result := c.Wrap(tt.color, tt.text)

			if tt.wantEqual {
				if result != tt.text {
					t.Errorf("Wrap(%q, %q) = %q, want %q", tt.color, tt.text, result, tt.text)
				}
				return
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("Wrap(%q, %q) = %q, want to contain %q", tt.color, tt.text, result, want)
				}
			}
		})
	}
}

func TestStrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "red ANSI code",
			input: "\033[31mERROR\033[0m",
			want:  "ERROR",
		},
		{
			name:  "yellow ANSI code",
			input: "\033[33mWARNING\033[0m",
			want:  "WARNING",
		},
		{
			name:  "multiple ANSI codes",
			input: "\033[31mRED\033[0m and \033[32mGREEN\033[0m",
			want:  "RED and GREEN",
		},
		{
			name:  "no ANSI codes",
			input: "plain text",
			want:  "plain text",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "complex ANSI sequence",
			input: "\033[1;31mBOLD RED\033[0m",
			want:  "BOLD RED",
		},
		{
			name:  "only ANSI codes",
			input: "\033[31m\033[0m",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Strip(tt.input)
			if result != tt.want {
				t.Errorf("Strip(%q) = %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}

func TestStripPreservesNonANSI(t *testing.T) {
	input := "Hello 世界 🌍"
	result := Strip(input)
	if result != input {
		t.Errorf("Strip(%q) = %q, want %q", input, result, input)
	}
}
