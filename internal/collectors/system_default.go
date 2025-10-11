package collectors

import (
	"context"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/host"
)

// DefaultSystemCollector gathers host metadata via gopsutil.
type DefaultSystemCollector struct{}

// NewSystemCollector returns the default implementation for all platforms.
func NewSystemCollector() SystemCollector {
	return DefaultSystemCollector{}
}

// CollectSystem implements SystemCollector.
func (DefaultSystemCollector) CollectSystem(ctx context.Context) (SystemInfo, error) {
	info, err := host.InfoWithContext(ctx)
	if err != nil {
		recordError("system", err)
		info = &host.InfoStat{Hostname: "unknown", Platform: runtime.GOOS}
	}
	var currentUser string
	homeDir := ""
	if u, err := user.Current(); err == nil {
		currentUser = u.Username
		homeDir = u.HomeDir
	}

	osName := prettyOS(info.Platform, info.PlatformFamily, info.OS)

	// Convert uptime safely from uint64 to int64 for time.Duration
	// Cap at max int64 to prevent overflow (292 years)
	uptime := info.Uptime
	if uptime > uint64(1<<63-1) {
		uptime = uint64(1<<63 - 1)
	}

	return SystemInfo{
		Hostname:    info.Hostname,
		OS:          osName,
		OSVersion:   info.PlatformVersion,
		Arch:        runtime.GOARCH,
		Uptime:      time.Duration(uptime) * time.Second, //nolint:gosec // G115: Overflow protected above (capped at max int64)
		CurrentUser: currentUser,
		HomeDir:     homeDir,
		Datetime:    time.Now(),
	}, nil
}

func prettyOS(platform, family, raw string) string {
	parts := []string{}
	if platform != "" {
		parts = append(parts, titleCase(platform))
	} else if raw != "" {
		parts = append(parts, titleCase(raw))
	}
	if family != "" && !strings.EqualFold(platform, family) {
		parts = append(parts, titleCase(family))
	}
	if len(parts) == 0 {
		return titleCase(raw)
	}
	return strings.Join(parts, " ")
}

// titleCase converts the first character of a string to uppercase.
// This is a simple ASCII-only replacement for deprecated strings.Title.
func titleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
