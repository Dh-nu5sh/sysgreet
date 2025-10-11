package banner

import (
	"fmt"
	"math"

	"github.com/veteranbv/sysgreet/internal/collectors"
	"github.com/veteranbv/sysgreet/internal/config"
)

// ResourceSectionBuilder renders memory/disk/CPU information.
type ResourceSectionBuilder struct{}

// Key returns resource identifier.
func (ResourceSectionBuilder) Key() string { return "resources" }

// Enabled reports whether any resource displays are enabled.
func (ResourceSectionBuilder) Enabled(cfg config.Config) bool {
	return cfg.Display.Memory || cfg.Display.Disk || cfg.Display.Load
}

// Build renders resource lines and attaches metadata for highlighting.
func (ResourceSectionBuilder) Build(snap collectors.Snapshot, cfg config.Config) (Section, bool) {
	var lines []string
	meta := map[string]any{}

	if cfg.Display.Memory && snap.Resources.Memory.Total > 0 {
		used := snap.Resources.Memory.Total - snap.Resources.Memory.Available
		pct := percent(used, snap.Resources.Memory.Total)
		lines = append(lines, fmt.Sprintf("Mem: %s free / %s (%d%% used)", humanBytes(snap.Resources.Memory.Available), humanBytes(snap.Resources.Memory.Total), pct))
		meta["memory_used_percent"] = pct
	}

	if cfg.Display.Disk && snap.Resources.Disk.Total > 0 {
		pct := percent(snap.Resources.Disk.Used, snap.Resources.Disk.Total)
		lines = append(lines, fmt.Sprintf("Disk: %s used / %s (%d%% used)", humanBytes(snap.Resources.Disk.Used), humanBytes(snap.Resources.Disk.Total), pct))
		meta["disk_used_percent"] = pct
	}

	if cfg.Display.Load {
		switch snap.Resources.CPU.Mode {
		case "usage":
			lines = append(lines, fmt.Sprintf("CPU: %.1f%%", snap.Resources.CPU.Usage))
			meta["cpu_usage_percent"] = snap.Resources.CPU.Usage
		case "load":
			lines = append(lines, fmt.Sprintf("CPU Load: %.2f %.2f %.2f", snap.Resources.CPU.Load1, snap.Resources.CPU.Load5, snap.Resources.CPU.Load15))
			meta["cpu_load_1"] = snap.Resources.CPU.Load1
		}
	}

	if len(lines) == 0 {
		return Section{}, false
	}
	return Section{Key: "resources", Title: "Resources", Lines: lines, Data: meta}, true
}

func humanBytes(value uint64) string {
	if value == 0 {
		return "0B"
	}
	const unit = 1024
	suffixes := []string{"B", "KB", "MB", "GB", "TB"}
	v := float64(value)
	i := 0
	for v >= unit && i < len(suffixes)-1 {
		v /= unit
		i++
	}
	return fmt.Sprintf("%.1f%s", v, suffixes[i])
}

func percent(part, total uint64) int {
	if total == 0 {
		return 0
	}
	return int(math.Round(float64(part) / float64(total) * 100))
}
