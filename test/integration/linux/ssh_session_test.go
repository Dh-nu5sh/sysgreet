package linux

import (
	"context"
	"testing"

	"github.com/veteranbv/hostinfo/internal/collectors"
)

func TestSessionCollectorUsesSSHConnection(t *testing.T) {
	t.Setenv("SSH_CONNECTION", "203.0.113.5 12345 10.0.0.5 22")
	t.Setenv("SSH_CLIENT", "198.51.100.7 54321 22")

	collector := collectors.NewSessionCollector()
	info, err := collector.CollectSession(context.Background())
	if err != nil {
		t.Fatalf("CollectSession error: %v", err)
	}
	if info.RemoteAddr != "203.0.113.5" {
		t.Fatalf("expected remote addr 203.0.113.5, got %s", info.RemoteAddr)
	}
	if info.Source != "SSH_CONNECTION" {
		t.Fatalf("expected source SSH_CONNECTION, got %s", info.Source)
	}
}
