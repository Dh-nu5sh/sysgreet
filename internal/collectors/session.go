package collectors

import (
	"context"
	"os"
	"strings"
)

// EnvSessionCollector reads SSH metadata from environment variables.
type EnvSessionCollector struct{}

// NewSessionCollector returns the default session collector.
func NewSessionCollector() SessionCollector {
	return EnvSessionCollector{}
}

// CollectSession implements SessionCollector.
func (EnvSessionCollector) CollectSession(ctx context.Context) (SessionInfo, error) {
	if addr := parseSSHEnv(os.Getenv("SSH_CONNECTION")); addr != "" {
		return SessionInfo{RemoteAddr: addr, Source: "SSH_CONNECTION"}, nil
	}
	if addr := parseSSHEnv(os.Getenv("SSH_CLIENT")); addr != "" {
		return SessionInfo{RemoteAddr: addr, Source: "SSH_CLIENT"}, nil
	}
	return SessionInfo{}, nil
}

func parseSSHEnv(value string) string {
	fields := strings.Fields(value)
	if len(fields) > 0 {
		return fields[0]
	}
	return ""
}
