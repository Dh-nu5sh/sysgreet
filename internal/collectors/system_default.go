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
	return SystemInfo{
		Hostname:    info.Hostname,
		OS:          osName,
		OSVersion:   info.PlatformVersion,
		Arch:        runtime.GOARCH,
		Uptime:      time.Duration(info.Uptime) * time.Second,
		CurrentUser: currentUser,
		HomeDir:     homeDir,
		Datetime:    time.Now(),
	}, nil
}

func prettyOS(platform, family, raw string) string {
	parts := []string{}
	if platform != "" {
		parts = append(parts, strings.Title(platform))
	} else if raw != "" {
		parts = append(parts, strings.Title(raw))
	}
	if family != "" && !strings.EqualFold(platform, family) {
		parts = append(parts, strings.Title(family))
	}
	if len(parts) == 0 {
		return strings.Title(raw)
	}
	return strings.Join(parts, " ")
}
