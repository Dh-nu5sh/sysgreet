package collectors

import (
	"context"
	"testing"
)

func TestResourceCollectorReturnsMetrics(t *testing.T) {
	collector := NewResourceCollector()
	info, err := collector.CollectResources(context.Background())
	if err != nil {
		t.Fatalf("CollectResources error: %v", err)
	}
	if info.Memory.Total == 0 {
		t.Fatalf("expected memory total > 0")
	}
	if info.Memory.Available == 0 {
		t.Fatalf("expected memory available > 0")
	}
	if info.Disk.Total == 0 {
		t.Fatalf("expected disk total > 0")
	}
	if info.Disk.Used == 0 {
		t.Fatalf("expected disk used > 0")
	}
	if info.CPU.Mode == "" {
		t.Fatalf("expected CPU mode set")
	}
}
