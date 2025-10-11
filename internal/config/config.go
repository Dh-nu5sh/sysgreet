package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

var (
	errUnsupportedFormat = errors.New("unsupported config format")
)

// Load returns the merged configuration, the path that was used, and an error if loading fails.
// Defaults are always applied; missing files are ignored.
func Load() (Config, string, error) {
	cfg := Default()
	candidatePaths := defaultConfigPaths()

	if custom := os.Getenv("HOSTINFO_CONFIG"); custom != "" {
		candidatePaths = append([]string{custom}, candidatePaths...)
	}

	var usedPath string
	for _, p := range candidatePaths {
		if p == "" {
			continue
		}
		expanded := expandPath(p)
		info, err := os.Stat(expanded)
		if err != nil {
			continue
		}
		if info.IsDir() {
			continue
		}

		data, err := os.ReadFile(expanded)
		if err != nil {
			return Config{}, "", fmt.Errorf("read config: %w", err)
		}

		var raw rawConfig
		switch strings.ToLower(filepath.Ext(expanded)) {
		case ".yaml", ".yml":
			if err := yaml.Unmarshal(data, &raw); err != nil {
				return Config{}, "", fmt.Errorf("parse yaml config: %w", err)
			}
		case ".toml":
			if err := toml.Unmarshal(data, &raw); err != nil {
				return Config{}, "", fmt.Errorf("parse toml config: %w", err)
			}
		default:
			return Config{}, "", fmt.Errorf("%w: %s", errUnsupportedFormat, expanded)
		}

		mergeConfig(&cfg, raw)
		usedPath = expanded
		break
	}

	applyEnvOverrides(&cfg)
	return cfg, usedPath, nil
}

func defaultConfigPaths() []string {
	home, _ := os.UserHomeDir()
	return []string{
		filepath.Join(home, ".config", "hostinfo", "config.yaml"),
		filepath.Join(home, ".config", "hostinfo", "config.yml"),
		filepath.Join(home, ".config", "hostinfo", "config.toml"),
		filepath.Join(home, ".hostinfo.yaml"),
		filepath.Join(home, ".hostinfo.yml"),
		filepath.Join(home, ".hostinfo.toml"),
	}
}

func expandPath(p string) string {
	if strings.HasPrefix(p, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, strings.TrimPrefix(p, "~"))
		}
	}
	return os.ExpandEnv(p)
}

func mergeConfig(base *Config, override rawConfig) {
	if override.Display != nil {
		if override.Display.Hostname != nil {
			base.Display.Hostname = *override.Display.Hostname
		}
		if override.Display.OS != nil {
			base.Display.OS = *override.Display.OS
		}
		if override.Display.IPAddresses != nil {
			base.Display.IPAddresses = *override.Display.IPAddresses
		}
		if override.Display.RemoteIP != nil {
			base.Display.RemoteIP = *override.Display.RemoteIP
		}
		if override.Display.Uptime != nil {
			base.Display.Uptime = *override.Display.Uptime
		}
		if override.Display.User != nil {
			base.Display.User = *override.Display.User
		}
		if override.Display.Memory != nil {
			base.Display.Memory = *override.Display.Memory
		}
		if override.Display.Disk != nil {
			base.Display.Disk = *override.Display.Disk
		}
		if override.Display.Load != nil {
			base.Display.Load = *override.Display.Load
		}
		if override.Display.Datetime != nil {
			base.Display.Datetime = *override.Display.Datetime
		}
		if override.Display.LastLogin != nil {
			base.Display.LastLogin = *override.Display.LastLogin
		}
	}
	if override.ASCII != nil {
		if override.ASCII.Font != nil && *override.ASCII.Font != "" {
			base.ASCII.Font = *override.ASCII.Font
		}
		if override.ASCII.Color != nil && *override.ASCII.Color != "" {
			base.ASCII.Color = *override.ASCII.Color
		}
		if override.ASCII.Monochrome != nil {
			base.ASCII.Monochrome = *override.ASCII.Monochrome
		}
	}
	if override.Layout != nil {
		if override.Layout.Sections != nil && len(*override.Layout.Sections) > 0 {
			base.Layout.Sections = append([]string{}, (*override.Layout.Sections)...)
		}
		if override.Layout.Compact != nil {
			base.Layout.Compact = *override.Layout.Compact
		}
	}
	if override.Network != nil {
		if override.Network.ShowInterfaceNames != nil {
			base.Network.ShowInterfaceNames = *override.Network.ShowInterfaceNames
		}
		if override.Network.MaxInterfaces != nil {
			base.Network.MaxInterfaces = *override.Network.MaxInterfaces
		}
	}
}

func applyEnvOverrides(cfg *Config) {
	if v, ok := lookupBool("HOSTINFO_DISPLAY_HOSTNAME"); ok {
		cfg.Display.Hostname = v
	}
	if v, ok := lookupBool("HOSTINFO_DISPLAY_OS"); ok {
		cfg.Display.OS = v
	}
	if v, ok := lookupBool("HOSTINFO_DISPLAY_IP_ADDRESSES"); ok {
		cfg.Display.IPAddresses = v
	}
	if v, ok := lookupBool("HOSTINFO_DISPLAY_REMOTE_IP"); ok {
		cfg.Display.RemoteIP = v
	}
	if v, ok := lookupBool("HOSTINFO_DISPLAY_UPTIME"); ok {
		cfg.Display.Uptime = v
	}
	if v, ok := lookupBool("HOSTINFO_DISPLAY_USER"); ok {
		cfg.Display.User = v
	}
	if v, ok := lookupBool("HOSTINFO_DISPLAY_MEMORY"); ok {
		cfg.Display.Memory = v
	}
	if v, ok := lookupBool("HOSTINFO_DISPLAY_DISK"); ok {
		cfg.Display.Disk = v
	}
	if v, ok := lookupBool("HOSTINFO_DISPLAY_LOAD"); ok {
		cfg.Display.Load = v
	}
	if v, ok := lookupBool("HOSTINFO_DISPLAY_DATETIME"); ok {
		cfg.Display.Datetime = v
	}
	if v, ok := lookupBool("HOSTINFO_DISPLAY_LAST_LOGIN"); ok {
		cfg.Display.LastLogin = v
	}

	if font := os.Getenv("HOSTINFO_ASCII_FONT"); font != "" {
		cfg.ASCII.Font = font
	}
	if color := os.Getenv("HOSTINFO_ASCII_COLOR"); color != "" {
		cfg.ASCII.Color = color
	}
	if v, ok := lookupBool("HOSTINFO_ASCII_MONOCHROME"); ok {
		cfg.ASCII.Monochrome = v
	}

	if v, ok := lookupBool("HOSTINFO_LAYOUT_COMPACT"); ok {
		cfg.Layout.Compact = v
	}
	if sections := os.Getenv("HOSTINFO_LAYOUT_SECTIONS"); sections != "" {
		parts := strings.Split(sections, ",")
		var cleaned []string
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				cleaned = append(cleaned, trimmed)
			}
		}
		if len(cleaned) > 0 {
			cfg.Layout.Sections = cleaned
		}
	}

	if v, ok := lookupBool("HOSTINFO_NETWORK_SHOW_INTERFACE_NAMES"); ok {
		cfg.Network.ShowInterfaceNames = v
	}
	if max := os.Getenv("HOSTINFO_NETWORK_MAX_INTERFACES"); max != "" {
		if parsed, err := parseInt(max); err == nil {
			cfg.Network.MaxInterfaces = parsed
		}
	}
}

func lookupBool(key string) (bool, bool) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return false, false
	}
	switch strings.ToLower(strings.TrimSpace(val)) {
	case "1", "true", "yes", "on":
		return true, true
	case "0", "false", "no", "off":
		return false, true
	default:
		return false, true
	}
}

func parseInt(input string) (int, error) {
	var value int
	_, err := fmt.Sscanf(strings.TrimSpace(input), "%d", &value)
	return value, err
}

type rawConfig struct {
	Display *rawDisplay `yaml:"display" toml:"display"`
	ASCII   *rawASCII   `yaml:"ascii" toml:"ascii"`
	Layout  *rawLayout  `yaml:"layout" toml:"layout"`
	Network *rawNetwork `yaml:"network" toml:"network"`
}

type rawDisplay struct {
	Hostname    *bool `yaml:"hostname" toml:"hostname"`
	OS          *bool `yaml:"os" toml:"os"`
	IPAddresses *bool `yaml:"ip_addresses" toml:"ip_addresses"`
	RemoteIP    *bool `yaml:"remote_ip" toml:"remote_ip"`
	Uptime      *bool `yaml:"uptime" toml:"uptime"`
	User        *bool `yaml:"user" toml:"user"`
	Memory      *bool `yaml:"memory" toml:"memory"`
	Disk        *bool `yaml:"disk" toml:"disk"`
	Load        *bool `yaml:"load" toml:"load"`
	Datetime    *bool `yaml:"datetime" toml:"datetime"`
	LastLogin   *bool `yaml:"last_login" toml:"last_login"`
}

type rawASCII struct {
	Font       *string `yaml:"font" toml:"font"`
	Color      *string `yaml:"color" toml:"color"`
	Monochrome *bool   `yaml:"monochrome" toml:"monochrome"`
}

type rawLayout struct {
	Compact  *bool     `yaml:"compact" toml:"compact"`
	Sections *[]string `yaml:"sections" toml:"sections"`
}

type rawNetwork struct {
	ShowInterfaceNames *bool `yaml:"show_interface_names" toml:"show_interface_names"`
	MaxInterfaces      *int  `yaml:"max_interfaces" toml:"max_interfaces"`
}
