package collectors

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

// DefaultResourceCollector gathers memory, disk, and CPU metrics.
type DefaultResourceCollector struct{}

// NewResourceCollector instantiates the default resource collector.
func NewResourceCollector() ResourceCollector {
	return DefaultResourceCollector{}
}

// CollectResources implements ResourceCollector with graceful degradation.
func (DefaultResourceCollector) CollectResources(ctx context.Context) (ResourceInfo, error) {
	var info ResourceInfo

	if vm, err := mem.VirtualMemoryWithContext(ctx); err == nil {
		info.Memory = MemoryInfo{Total: vm.Total, Available: vm.Available}
	} else {
		recordError("memory", err)
	}

	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		home = "."
	}
	if usage, err := disk.UsageWithContext(ctx, home); err == nil {
		info.Disk = DiskInfo{Total: usage.Total, Used: usage.Total - usage.Free}
	} else {
		recordError("disk", err)
	}

	if runtime.GOOS == "windows" {
		values, err := cpu.PercentWithContext(ctx, time.Millisecond*100, false)
		if err != nil {
			recordError("cpu", err)
		} else if len(values) > 0 {
			info.CPU = CPUInfo{Usage: values[0], Mode: "usage"}
		}
	} else {
		avg, err := load.AvgWithContext(ctx)
		if err != nil {
			recordError("cpu", err)
		} else {
			info.CPU = CPUInfo{Load1: avg.Load1, Load5: avg.Load5, Load15: avg.Load15, Mode: "load"}
		}
	}
	return info, nil
}
