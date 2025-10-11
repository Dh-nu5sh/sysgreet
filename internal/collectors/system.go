package collectors

import (
	"context"
	"time"
)

// SystemInfo captures host identity and session metadata.
type SystemInfo struct {
	Hostname    string
	OS          string
	OSVersion   string
	Arch        string
	Uptime      time.Duration
	CurrentUser string
	HomeDir     string
	Datetime    time.Time
}

// Address represents an IP address bound to an interface.
type Address struct {
	IP        string
	Interface string
}

// NetworkInfo summarizes primary and secondary addresses.
type NetworkInfo struct {
	Primary    *Address
	Additional []Address
}

// SessionInfo provides connection-specific metadata.
type SessionInfo struct {
	RemoteAddr string
	Source     string
}

// MemoryInfo captures memory usage snapshot.
type MemoryInfo struct {
	Total     uint64
	Available uint64
}

// DiskInfo captures disk usage snapshot.
type DiskInfo struct {
	Total uint64
	Used  uint64
}

// CPUInfo captures load or usage metrics.
type CPUInfo struct {
	Load1  float64
	Load5  float64
	Load15 float64
	Usage  float64
	Mode   string // "load" or "usage"
}

// ResourceInfo bundles resource metrics for display.
type ResourceInfo struct {
	Memory MemoryInfo
	Disk   DiskInfo
	CPU    CPUInfo
}

// LastLoginInfo contains last successful login data.
type LastLoginInfo struct {
	Timestamp time.Time
	Source    string
}

// Snapshot aggregates all collector outputs for banner rendering.
type Snapshot struct {
	System    SystemInfo
	Network   NetworkInfo
	Session   SessionInfo
	Resources ResourceInfo
	LastLogin *LastLoginInfo
}

// SystemCollector defines host identity collection behaviour.
type SystemCollector interface {
	CollectSystem(ctx context.Context) (SystemInfo, error)
}

// NetworkCollector defines network snapshot behaviour.
type NetworkCollector interface {
	CollectNetwork(ctx context.Context) (NetworkInfo, error)
}

// ResourceCollector defines resource metrics behaviour.
type ResourceCollector interface {
	CollectResources(ctx context.Context) (ResourceInfo, error)
}

// SessionCollector defines remote session detection behaviour.
type SessionCollector interface {
	CollectSession(ctx context.Context) (SessionInfo, error)
}

// LastLoginCollector defines retrieval of last-login metadata.
type LastLoginCollector interface {
	CollectLastLogin(ctx context.Context) (*LastLoginInfo, error)
}

// Providers groups all collectors used to build a snapshot.
type Providers struct {
	System    SystemCollector
	Network   NetworkCollector
	Resources ResourceCollector
	Session   SessionCollector
	LastLogin LastLoginCollector
}

// Gather builds a Snapshot using the configured providers. Missing collectors are tolerated to allow graceful degradation.
func (p Providers) Gather(ctx context.Context) Snapshot {
	var snap Snapshot
	if p.System != nil {
		if sys, err := p.System.CollectSystem(ctx); err == nil {
			fmtSystem(&snap.System, sys)
		} else {
			recordError("system", err)
		}
	}
	if p.Network != nil {
		if netInfo, err := p.Network.CollectNetwork(ctx); err == nil {
			snap.Network = netInfo
		} else {
			recordError("network", err)
		}
	}
	if p.Resources != nil {
		if res, err := p.Resources.CollectResources(ctx); err == nil {
			snap.Resources = res
		} else {
			recordError("resources", err)
		}
	}
	if p.Session != nil {
		if session, err := p.Session.CollectSession(ctx); err == nil {
			snap.Session = session
		} else {
			recordError("session", err)
		}
	}
	if p.LastLogin != nil {
		if last, err := p.LastLogin.CollectLastLogin(ctx); err == nil {
			snap.LastLogin = last
		} else {
			recordError("last_login", err)
		}
	}
	return snap
}

func fmtSystem(dst *SystemInfo, src SystemInfo) {
	*dst = src
}
