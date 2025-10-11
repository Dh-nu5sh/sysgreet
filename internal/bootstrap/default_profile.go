package bootstrap

import (
	"time"

	"gopkg.in/yaml.v3"

	"github.com/veteranbv/sysgreet/internal/config"
)

// renderDefaultConfig marshals the default sysgreet configuration into YAML with metadata.
func renderDefaultConfig(now time.Time) ([]byte, error) {
	cfg := config.Default()
	cfg.Version = config.SchemaVersion
	cfg.CreatedAt = now.UTC().Format(time.RFC3339)
	return yaml.Marshal(cfg)
}
