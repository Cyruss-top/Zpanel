package monitor

// Overview 系统概览
type Overview struct {
	CPU     CPUInfo     `json:"cpu"`
	Memory  MemoryInfo  `json:"memory"`
	Disk    DiskInfo    `json:"disk"`
	Load    LoadInfo    `json:"load"`
	Uptime  int64       `json:"uptime_seconds"`
	OS      string      `json:"os"`
	Arch    string      `json:"arch"`
	Platform string     `json:"platform"`
}

type CPUInfo struct {
	UsagePercent float64 `json:"usage_percent"`
	Cores        int     `json:"cores"`
}

type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Available   uint64  `json:"available"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Available   uint64  `json:"available"`
	UsedPercent float64 `json:"used_percent"`
	Path        string  `json:"path"`
}

type LoadInfo struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// Collect 采集系统指标（Linux 读 /proc，其他平台降级）
func Collect(dataPath string) (*Overview, error) {
	return collect(dataPath)
}
