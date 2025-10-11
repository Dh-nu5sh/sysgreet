//go:build darwin

package darwin

import (
	"context"
	"math"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/veteranbv/hostinfo/internal/collectors"
)

func TestResourceCollectorMatchesSystemStats(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("macOS-specific accuracy test")
	}
	ctx := context.Background()
	collector := collectors.NewResourceCollector()
	info, err := collector.CollectResources(ctx)
	if err != nil {
		t.Fatalf("CollectResources error: %v", err)
	}

	vm, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		t.Fatalf("VirtualMemory error: %v", err)
	}
	if !withinTolerance(float64(info.Memory.Total), float64(vm.Total), 0.05) {
		t.Fatalf("memory total mismatch: got %d expected approx %d", info.Memory.Total, vm.Total)
	}

	homeUsage, err := disk.UsageWithContext(ctx, filepath.Clean("."))
	if err != nil {
		t.Fatalf("Disk usage error: %v", err)
	}
	if !withinTolerance(float64(info.Disk.Total), float64(homeUsage.Total), 0.10) {
		t.Fatalf("disk total mismatch got %d expected approx %d", info.Disk.Total, homeUsage.Total)
	}
}

func withinTolerance(observed, expected, tolerance float64) bool {
	if expected == 0 {
		return false
	}
	return math.Abs(observed-expected)/expected <= tolerance
}
