package render

import (
	"strings"
	"testing"

	"github.com/veteranbv/hostinfo/internal/banner"
	"github.com/veteranbv/hostinfo/internal/config"
)

func TestRenderer_Render(t *testing.T) {
	tests := []struct {
		name         string
		output       banner.Output
		cfg          config.Config
		wantContains []string
	}{
		{
			name: "basic header and sections",
			output: banner.Output{
				Header: banner.Header{
					Art:   "ASCII ART",
					Lines: []string{"Linux 6.0 (x86_64)"},
				},
				Sections: []banner.Section{
					{
						Key:   "system",
						Title: "System",
						Lines: []string{"Uptime: 1d 2h 30m"},
					},
				},
			},
			cfg:          config.Default(),
			wantContains: []string{"ASCII ART", "Linux 6.0 (x86_64)", "System", "Uptime: 1d 2h 30m"},
		},
		{
			name: "multiple sections with ordering",
			output: banner.Output{
				Header: banner.Header{
					Art: "HOST",
				},
				Sections: []banner.Section{
					{
						Key:   "resources",
						Title: "Resources",
						Lines: []string{"Mem: 8GB"},
					},
					{
						Key:   "system",
						Title: "System",
						Lines: []string{"Uptime: 1d"},
					},
				},
			},
			cfg: config.Config{
				Layout: config.LayoutConfig{
					Sections: []string{"system", "resources"},
				},
			},
			wantContains: []string{"System", "Resources"},
		},
		{
			name: "empty sections are skipped",
			output: banner.Output{
				Header: banner.Header{
					Art: "HOST",
				},
				Sections: []banner.Section{
					{
						Key:   "system",
						Title: "System",
						Lines: []string{},
					},
					{
						Key:   "network",
						Title: "Network",
						Lines: []string{"Primary: 192.168.1.1"},
					},
				},
			},
			cfg:          config.Default(),
			wantContains: []string{"Network", "Primary: 192.168.1.1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRenderer(true) // Disable color for deterministic output
			result := r.Render(tt.output, tt.cfg)

			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("Render() result missing %q\nGot:\n%s", want, result)
				}
			}
		})
	}
}

func TestRenderer_RenderCompact(t *testing.T) {
	tests := []struct {
		name         string
		output       banner.Output
		cfg          config.Config
		wantContains []string
		separator    string
	}{
		{
			name: "compact mode with separator",
			output: banner.Output{
				Header: banner.Header{
					Art:   "HOST",
					Lines: []string{"Linux 6.0"},
				},
				Sections: []banner.Section{
					{
						Key:   "system",
						Title: "System",
						Lines: []string{"Uptime: 1d"},
					},
				},
			},
			cfg: config.Config{
				Layout: config.LayoutConfig{
					Compact: true,
				},
			},
			wantContains: []string{"HOST", "Linux 6.0", "System", "Uptime: 1d"},
			separator:    " | ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRenderer(true)
			result := r.Render(tt.output, tt.cfg)

			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("Render() compact result missing %q\nGot: %s", want, result)
				}
			}

			if !strings.Contains(result, tt.separator) {
				t.Errorf("Render() compact result missing separator %q", tt.separator)
			}
		})
	}
}

func TestOrderSections(t *testing.T) {
	tests := []struct {
		name     string
		sections []banner.Section
		desired  []string
		wantKeys []string
	}{
		{
			name: "ordered by desired list",
			sections: []banner.Section{
				{Key: "resources"},
				{Key: "system"},
				{Key: "network"},
			},
			desired:  []string{"system", "network", "resources"},
			wantKeys: []string{"system", "network", "resources"},
		},
		{
			name: "partial ordering",
			sections: []banner.Section{
				{Key: "resources"},
				{Key: "system"},
				{Key: "network"},
			},
			desired:  []string{"system"},
			wantKeys: []string{"system"},
		},
		{
			name: "alphabetical fallback when desired is empty",
			sections: []banner.Section{
				{Key: "resources"},
				{Key: "system"},
				{Key: "network"},
			},
			desired:  []string{},
			wantKeys: []string{"network", "resources", "system"},
		},
		{
			name: "desired order with missing keys",
			sections: []banner.Section{
				{Key: "system"},
			},
			desired:  []string{"network", "system", "resources"},
			wantKeys: []string{"system"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := orderSections(tt.sections, tt.desired)

			if len(result) != len(tt.wantKeys) {
				t.Fatalf("got %d sections, want %d", len(result), len(tt.wantKeys))
			}

			for i, want := range tt.wantKeys {
				if result[i].Key != want {
					t.Errorf("section[%d].Key = %q, want %q", i, result[i].Key, want)
				}
			}
		})
	}
}

func TestRenderer_HighlightResource(t *testing.T) {
	tests := []struct {
		name         string
		section      banner.Section
		line         string
		disableColor bool
		wantContains string
	}{
		{
			name: "memory high usage (red)",
			section: banner.Section{
				Key: "resources",
				Data: map[string]any{
					"memory_used_percent": 95,
				},
			},
			line:         "Mem: 15GB used / 16GB",
			disableColor: false,
			wantContains: "Mem:",
		},
		{
			name: "memory warning usage (yellow)",
			section: banner.Section{
				Key: "resources",
				Data: map[string]any{
					"memory_used_percent": 80,
				},
			},
			line:         "Mem: 12GB used / 16GB",
			disableColor: false,
			wantContains: "Mem:",
		},
		{
			name: "disk normal usage (no color)",
			section: banner.Section{
				Key: "resources",
				Data: map[string]any{
					"disk_used_percent": 50,
				},
			},
			line:         "Disk: 250GB used / 500GB",
			disableColor: false,
			wantContains: "Disk:",
		},
		{
			name: "no data map",
			section: banner.Section{
				Key:  "resources",
				Data: nil,
			},
			line:         "Mem: 8GB used / 16GB",
			disableColor: false,
			wantContains: "Mem:",
		},
		{
			name: "color disabled",
			section: banner.Section{
				Key: "resources",
				Data: map[string]any{
					"memory_used_percent": 95,
				},
			},
			line:         "Mem: 15GB used / 16GB",
			disableColor: true,
			wantContains: "Mem: 15GB used / 16GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRenderer(tt.disableColor)
			result := r.highlightResource(tt.section, tt.line)

			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("highlightResource() = %q, want to contain %q", result, tt.wantContains)
			}

			// When color is disabled, output should equal input
			if tt.disableColor && result != tt.line {
				t.Errorf("highlightResource() with disabled color = %q, want %q", result, tt.line)
			}
		})
	}
}

func TestRenderer_WrapForPercent(t *testing.T) {
	tests := []struct {
		name         string
		pct          int
		line         string
		disableColor bool
		wantColor    bool
	}{
		{"critical threshold", 95, "Mem: 95%", false, true},
		{"warning threshold", 80, "Mem: 80%", false, true},
		{"normal threshold", 50, "Mem: 50%", false, false},
		{"zero percent", 0, "Mem: 0%", false, false},
		{"color disabled", 95, "Mem: 95%", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRenderer(tt.disableColor)
			result := r.wrapForPercent(tt.pct, tt.line)

			hasColor := result != tt.line
			if hasColor != tt.wantColor {
				t.Errorf("wrapForPercent(%d, %q) hasColor=%v, want %v", tt.pct, tt.line, hasColor, tt.wantColor)
			}
		})
	}
}
