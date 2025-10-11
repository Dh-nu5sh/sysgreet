package ascii

import "strings"

// RenderOptions describes how ASCII art should be generated.
type RenderOptions struct {
	Font       string
	Color      string
	Monochrome bool
	Uppercase  bool
}

// RenderHostname renders the hostname into ASCII art using the configured options.
func (r *Renderer) RenderHostname(hostname string, opts RenderOptions) (string, string, string, error) {
	text := strings.TrimSpace(hostname)
	if text == "" {
		text = "hostinfo"
	}
	if opts.Uppercase {
		text = strings.ToUpper(text)
	}
	return r.Render(text, opts.Font, opts.Color, opts.Monochrome)
}
