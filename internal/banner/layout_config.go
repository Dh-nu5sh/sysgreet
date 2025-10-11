package banner

import "github.com/veteranbv/sysgreet/internal/config"

// BuildersForConfig returns section builders ordered according to configuration.
func BuildersForConfig(cfg config.Config) []Builder {
	available := map[string]Builder{
		"system":    SystemSectionBuilder{},
		"network":   NetworkSectionBuilder{},
		"resources": ResourceSectionBuilder{},
	}
	var builders []Builder
	for _, key := range cfg.Layout.Sections {
		if b, ok := available[key]; ok {
			builders = append(builders, b)
		}
	}
	if len(builders) == 0 {
		builders = []Builder{SystemSectionBuilder{}, NetworkSectionBuilder{}, ResourceSectionBuilder{}}
	}
	return builders
}
