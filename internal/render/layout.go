package render

import (
	"sort"
	"strings"

	"github.com/veteranbv/hostinfo/internal/banner"
	"github.com/veteranbv/hostinfo/internal/config"
)

// Renderer formats banner output into terminal-friendly text.
type Renderer struct {
	colorizer Colorizer
}

// NewRenderer instantiates a renderer with optional color suppression.
func NewRenderer(disableColor bool) Renderer {
	return Renderer{colorizer: NewColorizer(disableColor)}
}

// Render produces the final banner string.
func (r Renderer) Render(out banner.Output, cfg config.Config) string {
	if cfg.Layout.Compact {
		return r.renderCompact(out, cfg)
	}

	var builder strings.Builder
	builder.WriteString(out.Header.Art)
	if len(out.Header.Lines) > 0 {
		builder.WriteString("\n")
		for _, line := range out.Header.Lines {
			builder.WriteString(line)
			builder.WriteString("\n")
		}
	}

	sections := orderSections(out.Sections, cfg.Layout.Sections)
	for _, section := range sections {
		if len(section.Lines) == 0 {
			continue
		}
		builder.WriteString("\n")
		builder.WriteString(section.Title)
		builder.WriteString("\n")
		for _, line := range section.Lines {
			formatted := line
			if section.Key == "resources" {
				formatted = r.highlightResource(section, line)
			}
			builder.WriteString("  ")
			builder.WriteString(formatted)
			builder.WriteString("\n")
		}
	}
	return strings.TrimRight(builder.String(), "\n")
}

func (r Renderer) renderCompact(out banner.Output, cfg config.Config) string {
	parts := []string{Strip(out.Header.Art)}
	parts = append(parts, out.Header.Lines...)
	sections := orderSections(out.Sections, cfg.Layout.Sections)
	for _, section := range sections {
		if len(section.Lines) == 0 {
			continue
		}
		parts = append(parts, section.Title)
		parts = append(parts, section.Lines...)
	}
	return strings.Join(parts, " | ")
}

func orderSections(sections []banner.Section, desired []string) []banner.Section {
	lookup := make(map[string]banner.Section)
	keys := []string{}
	for _, s := range sections {
		lookup[s.Key] = s
		keys = append(keys, s.Key)
	}
	var ordered []banner.Section
	for _, key := range desired {
		if sec, ok := lookup[key]; ok {
			ordered = append(ordered, sec)
		}
	}
	if len(ordered) == 0 {
		sort.Strings(keys)
		for _, key := range keys {
			ordered = append(ordered, lookup[key])
		}
	}
	return ordered
}

func (r Renderer) highlightResource(section banner.Section, line string) string {
	data := section.Data
	if data == nil {
		return line
	}
	switch {
	case strings.HasPrefix(line, "Mem:"):
		if pct, ok := data["memory_used_percent"].(int); ok {
			return r.wrapForPercent(pct, line)
		}
	case strings.HasPrefix(line, "Disk:"):
		if pct, ok := data["disk_used_percent"].(int); ok {
			return r.wrapForPercent(pct, line)
		}
	case strings.HasPrefix(line, "CPU:") && strings.Contains(line, "%"):
		if pct, ok := data["cpu_usage_percent"].(float64); ok {
			return r.wrapForPercent(int(pct+0.5), line)
		}
	}
	return line
}

func (r Renderer) wrapForPercent(pct int, line string) string {
	color := ""
	switch {
	case pct >= 90:
		color = "red"
	case pct >= 75:
		color = "yellow"
	case pct >= 0:
		color = "green"
	}
	if color == "" || color == "green" {
		return line
	}
	return r.colorizer.Wrap(color, line)
}
