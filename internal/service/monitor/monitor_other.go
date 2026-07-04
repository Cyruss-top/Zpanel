//go:build !linux

package monitor

import (
	"runtime"
	"time"
)

func collect(dataPath string) (*Overview, error) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	total := ms.Sys
	used := ms.Alloc
	avail := total - used
	if total < used {
		total = ms.HeapSys
		used = ms.HeapAlloc
	}
	pct := 0.0
	if total > 0 {
		pct = float64(used) / float64(total) * 100
	}

	return &Overview{
		CPU: CPUInfo{
			UsagePercent: 0,
			Cores:        runtime.NumCPU(),
		},
		Memory: MemoryInfo{
			Total:       total,
			Used:        used,
			Available:   avail,
			UsedPercent: pct,
		},
		Disk: DiskInfo{
			Path: dataPath,
		},
		Load:     LoadInfo{},
		Uptime:   int64(time.Since(startedAt).Seconds()),
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Platform: runtime.GOOS + " (limited metrics)",
	}, nil
}

var startedAt = time.Now()
