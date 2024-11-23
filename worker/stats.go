package worker

import (
	"log"

	"github.com/c9s/goprocinfo/linux"
)

type Stats struct {
	MemStats *linux.MemInfo
	DiskStats *linux.Disk
	CpuStats *linux.CPUStat
	LoadStats *linux.LoadAvg
	TaskCount int
}
func (s *Stats) MemTotalKb() uint64 {
	return s.MemStats.MemTotal
}
func (s *Stats) MemAvailableKb() uint64 {
	return s.MemStats.MemAvailable
}
func (s *Stats) MemUsedKb() uint64 {
	return s.MemStats.MemTotal - s.MemStats.MemAvailable
}
func (s *Stats) MemUsedPercent() uint64 {
	return uint64((float64(s.MemUsedKb()) / float64(s.MemTotalKb())) * 100)
}

func (s *Stats) DiskTotal() uint64 {
	return s.DiskStats.All
}
func (s *Stats) DiskFree() uint64 {
	return s.DiskStats.Free
}
func (s *Stats) DiskUsed() uint64 {
	return s.DiskStats.Used
}

/*
$ cat /proc/loadavg

host@ubuntu:~$ cat /proc/loadavg
0.37 0.46 0.47 1/1950 69695
host@ubuntu:~$ uptime
 23:19:50 up 20:53,  2 users,  load average: 0.39, 0.46, 0.47

 https://pkg.go.dev/github.com/c9s/goprocinfo/linux#CPUStat
type CPUStat struct {
    Id        string `json:"id"`
    User      uint64 `json:"user"`
    Nice      uint64 `json:"nice"`
    System    uint64 `json:"system"`
    Idle      uint64 `json:"idle"`
    IOWait    uint64 `json:"iowait"`
    IRQ       uint64 `json:"irq"`
    SoftIRQ   uint64 `json:"softirq"`
    Steal     uint64 `json:"steal"`
    Guest     uint64 `json:"guest"`
    GuestNice uint64 `json:"guest_nice"`
}



*/

func (s *Stats) CpuUsage() float64{
	// sum of values for idle states
	idle := s.CpuStats.Idle + s.CpuStats.IOWait
	// sum for values for non-idle states
	nonIdle := s.CpuStats.User + s.CpuStats.Nice + s.CpuStats.System +
	s.CpuStats.IRQ + s.CpuStats.SoftIRQ + s.CpuStats.Steal
	// total sum of idle and non-idle 
	total := idle + nonIdle
	

	if total == 0 {
		return 0.00
	}
	// total cpu usage = (total - idle) / total
	// https://stackoverflow.com/questions/23367857/accurate-calculation-of-cpu-usage-given-in-percentage-in-linux
	return (float64(total) - float64(idle)) / float64(total)
}

func GetStats() *Stats {
	return &Stats{
		MemStats: GetMemoryInfo(),
		DiskStats: GetDiskInfo(),
		CpuStats: GetCpuStats(),
		LoadStats: GetLoadAvg(),
	}
}

func GetMemoryInfo() *linux.MemInfo {
	memstats, err := linux.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Printf("error reading /proc/meminfo")
		return &linux.MemInfo{}
	}
	return memstats
}
// https://godoc.org/github.com/c9s/goprocinfo/linux#Disk
func GetDiskInfo() *linux.Disk {
	diskstats, err := linux.ReadDisk("/")
	if err != nil {
		log.Printf("error reading from /")
		return &linux.Disk{}
	}
	return diskstats
}
// https://godoc.org/github.com/c9s/goprocinfo/linux#CPUStat
func GetCpuStats() *linux.CPUStat {
	stats, err := linux.ReadStat("/proc/stat")
	if err != nil {
		log.Printf("error reading form /proc/stat")
		return &linux.CPUStat{}
	}
	return &stats.CPUStatAll
}
//https://godoc.org/github.com/c9s/goprocinfo/linux#LoadAvg
func GetLoadAvg() *linux.LoadAvg {
	loadavg, err := linux.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		log.Printf("error reading form /proc/loadavg")
		return &linux.LoadAvg{}
	}

	return loadavg
}