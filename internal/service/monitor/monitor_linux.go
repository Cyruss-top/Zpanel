//go:build linux

package monitor

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func collect(dataPath string) (*Overview, error) {
	mem, err := readMemInfo()
	if err != nil {
		return nil, err
	}
	cpuUsage, cores, err := readCPUUsage()
	if err != nil {
		cpuUsage = 0
		cores = runtime.NumCPU()
	}
	disk, err := readDisk(dataPath)
	if err != nil {
		disk = DiskInfo{Path: dataPath}
	}
	load, _ := readLoadAvg()
	uptime, _ := readUptime()

	return &Overview{
		CPU:      CPUInfo{UsagePercent: cpuUsage, Cores: cores},
		Memory:   mem,
		Disk:     disk,
		Load:     load,
		Uptime:   uptime,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Platform: "linux",
	}, nil
}

func readMemInfo() (MemoryInfo, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemoryInfo{}, err
	}
	defer f.Close()

	vals := map[string]uint64{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		parts := strings.Fields(sc.Text())
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSuffix(parts[0], ":")
		v, _ := strconv.ParseUint(parts[1], 10, 64)
		vals[key] = v * 1024
	}
	total := vals["MemTotal"]
	avail := vals["MemAvailable"]
	if avail == 0 {
		avail = vals["MemFree"] + vals["Buffers"] + vals["Cached"]
	}
	used := total - avail
	pct := 0.0
	if total > 0 {
		pct = float64(used) / float64(total) * 100
	}
	return MemoryInfo{Total: total, Used: used, Available: avail, UsedPercent: pct}, nil
}

var lastCPU [2]struct{ total, idle uint64 }

func readCPUUsage() (float64, int, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, runtime.NumCPU(), err
	}
	line := strings.Split(string(data), "\n")[0]
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0, runtime.NumCPU(), fmt.Errorf("bad /proc/stat")
	}
	var total, idle uint64
	for i := 1; i < len(fields); i++ {
		v, _ := strconv.ParseUint(fields[i], 10, 64)
		total += v
		if i == 4 {
			idle = v
		}
	}
	cores := runtime.NumCPU()
	idx := 0
	if lastCPU[0].total > 0 {
		idx = 1
	}
	prev := lastCPU[1-idx]
	cur := struct{ total, idle uint64 }{total, idle}
	lastCPU[idx] = cur

	if prev.total == 0 {
		time.Sleep(200 * time.Millisecond)
		return readCPUUsage()
	}
	dTotal := float64(cur.total - prev.total)
	dIdle := float64(cur.idle - prev.idle)
	if dTotal == 0 {
		return 0, cores, nil
	}
	return (1 - dIdle/dTotal) * 100, cores, nil
}

func readDisk(path string) (DiskInfo, error) {
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return DiskInfo{}, err
	}
	total := st.Blocks * uint64(st.Bsize)
	avail := st.Bavail * uint64(st.Bsize)
	used := total - avail
	pct := 0.0
	if total > 0 {
		pct = float64(used) / float64(total) * 100
	}
	return DiskInfo{Total: total, Used: used, Available: avail, UsedPercent: pct, Path: path}, nil
}

func readLoadAvg() (LoadInfo, error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return LoadInfo{}, err
	}
	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return LoadInfo{}, fmt.Errorf("bad loadavg")
	}
	l1, _ := strconv.ParseFloat(fields[0], 64)
	l5, _ := strconv.ParseFloat(fields[1], 64)
	l15, _ := strconv.ParseFloat(fields[2], 64)
	return LoadInfo{Load1: l1, Load5: l5, Load15: l15}, nil
}

func readUptime() (int64, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return 0, fmt.Errorf("bad uptime")
	}
	f, _ := strconv.ParseFloat(fields[0], 64)
	return int64(f), nil
}
