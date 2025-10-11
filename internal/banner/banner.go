package banner

import (
	"context"
	"errors"
	"strings"

	"github.com/veteranbv/hostinfo/internal/ascii"
	"github.com/veteranbv/hostinfo/internal/collectors"
	"github.com/veteranbv/hostinfo/internal/config"
)

// Section represents a rendered section of the banner body.
type Section struct {
	Key   string
	Title string
	Lines []string
	Data  map[string]any
}

// Header captures ASCII art metadata for the hostname banner.
type Header struct {
	Art   string
	Font  string
	Color string
	Lines []string
}

// Output encapsulates the full banner content ready for layout rendering.
type Output struct {
	Header   Header
	Sections []Section
}

// Builder generates a banner section when enabled.
type Builder interface {
	Key() string
	Enabled(cfg config.Config) bool
	Build(snap collectors.Snapshot, cfg config.Config) (Section, bool)
}

// Banner orchestrates collectors and builders to produce the final output.
type Banner struct {
	providers collectors.Providers
	ascii     *ascii.Renderer
	builders  []Builder
}

// New creates a Banner orchestrator.
func New(providers collectors.Providers, renderer *ascii.Renderer, builders []Builder) (*Banner, error) {
	if renderer == nil {
		return nil, errors.New("ascii renderer is required")
	}
	return &Banner{providers: providers, ascii: renderer, builders: builders}, nil
}

// Build produces the banner output using the provided configuration.
func (b *Banner) Build(ctx context.Context, cfg config.Config) (Output, collectors.Snapshot, error) {
	snap := b.providers.Gather(ctx)
	header := b.buildHeader(snap, cfg)
	sections := b.buildSections(snap, cfg)
	return Output{Header: header, Sections: sections}, snap, nil
}

func (b *Banner) buildHeader(snap collectors.Snapshot, cfg config.Config) Header {
	name := snap.System.Hostname
	if strings.TrimSpace(name) == "" {
		name = "hostinfo"
	}
	art, font, color, err := b.ascii.RenderHostname(name, ascii.RenderOptions{
		Font:       cfg.ASCII.Font,
		Color:      cfg.ASCII.Color,
		Monochrome: cfg.ASCII.Monochrome,
		Uppercase:  true,
	})
	if err != nil {
		// Fallback to plain text when ASCII rendering fails.
		art = name
		font = "plain"
		color = "reset"
	}
	lines := []string{}
	if cfg.Display.OS {
		line := snap.System.OS
		if snap.System.OSVersion != "" {
			line += " " + snap.System.OSVersion
		}
		if snap.System.Arch != "" {
			line += " (" + snap.System.Arch + ")"
		}
		lines = append(lines, line)
	}
	return Header{Art: art, Font: font, Color: color, Lines: lines}
}

func (b *Banner) buildSections(snap collectors.Snapshot, cfg config.Config) []Section {
	var sections []Section
	for _, builder := range b.builders {
		if builder == nil {
			continue
		}
		if !builder.Enabled(cfg) {
			continue
		}
		if section, ok := builder.Build(snap, cfg); ok && len(section.Lines) > 0 {
			sections = append(sections, section)
		}
	}
	return sections
}
