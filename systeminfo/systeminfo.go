package systeminfo

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

// Get general CPU Info
func GetCPUinfo() (cpu.InfoStat, error) {

	cpuInfo, err := cpu.Info()
	if err != nil {
		return cpu.InfoStat{}, fmt.Errorf("unable to get CPU info: %w", err)
	}

	return cpuInfo[0], nil
}

// Returns percentual value of total CPU usage
func GetCPUPercent() (float64, error) {
	cpuPercentage, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		return 0.0, fmt.Errorf("unable to get CPU usage percent: %w", err)
	}

	return cpuPercentage[0], nil
}

// Return Cpu Loads
func GetCPULoad() ([]cpu.TimesStat, error) {

	cpuLoad, err := cpu.Times(false)
	if err != nil {
		return []cpu.TimesStat{}, fmt.Errorf("unable to get CPU times: %w", err)
	}

	currLoad := cpuLoad[0]

	totalLoad := currLoad.Guest + currLoad.Idle + currLoad.Iowait + currLoad.Irq +
		currLoad.Nice + currLoad.Softirq + currLoad.Steal + currLoad.System + currLoad.User

	// Convert loads to percentual values
	currLoad.Guest = (currLoad.Guest / totalLoad) * 100
	currLoad.Idle = (currLoad.Idle / totalLoad) * 100
	currLoad.Iowait = (currLoad.Iowait / totalLoad) * 100
	currLoad.Irq = (currLoad.Irq / totalLoad) * 100
	currLoad.Nice = (currLoad.Nice / totalLoad) * 100
	currLoad.Softirq = (currLoad.Softirq / totalLoad) * 100
	currLoad.Steal = (currLoad.Steal / totalLoad) * 100
	currLoad.System = (currLoad.System / totalLoad) * 100
	currLoad.User = (currLoad.User / totalLoad) * 100

	cpuLoad[0] = currLoad
	return cpuLoad, nil

}

// Returns Memory usage statistics
func GetMEMLoad() (*mem.VirtualMemoryStat, error) {

	v, err := mem.VirtualMemory()
	if err != nil {
		return &mem.VirtualMemoryStat{}, fmt.Errorf("unable to get memory state: %w", err)
	}

	return &mem.VirtualMemoryStat{
		Total:       v.Total,
		Available:   v.Available,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
		Free:        v.Free,
		Buffers:     v.Buffers,
		Cached:      v.Cached,
	}, nil

}

type DiskInfo struct {
	Partition disk.PartitionStat
	Fstype    string
	Total     uint64
	Free      uint64
	Used      uint64
}

func GetDISKUse() ([]DiskInfo, error) {

	partitions, err := disk.Partitions(true)
	var disks []DiskInfo
	if err != nil {
		return disks, fmt.Errorf("unable to get disk info: %w", err)
	}

	for _, p := range partitions {
		diskInfo := DiskInfo{}
		diskInfo.Partition = p
		usageStat, err := disk.Usage(p.Mountpoint)
		if err != nil {
			return disks, fmt.Errorf("unable to get disk stats: %w", err)
		}
		diskInfo.Free = usageStat.Free
		diskInfo.Fstype = usageStat.Fstype
		diskInfo.Total = usageStat.Total
		diskInfo.Used = usageStat.Used

		disks = append(disks, diskInfo)
	}

	if len(disks) == 0 {
		return disks, errors.New("couldn't get disk partitions' stats")
	}

	sort.Slice(disks, func(i, j int) bool {
		return disks[i].Total > disks[j].Total
	})

	return disks, nil
}

type ProcessInfo struct {
	PID     int32
	Name    string
	CPU     float64
	Memory  uint64
	Runtime string
	Status  []string
}

func GetProcessInfo(n int) ([]ProcessInfo, error) {

	processes, err := process.Processes()
	if err != nil {
		return []ProcessInfo{}, fmt.Errorf("unable to read running processes: %w", err)
	}

	var processesInfo []ProcessInfo

	for _, p := range processes {
		proc := ProcessInfo{}
		proc.PID = p.Pid

		proc.Name, err = p.Name()
		if err != nil {
			proc.Name = "N/A"
		}

		proc.Status, err = p.Status()
		if err != nil {
			proc.Status = []string{"Unknown"}
		}

		memoryInfo, err := p.MemoryInfo()

		// If for loop is not broken after a memoryInfo error
		// a runtime error occurs
		if err != nil {
			processesInfo = append(processesInfo, ProcessInfo{
				PID:    proc.PID,
				Name:   proc.Name,
				Status: proc.Status,
				Memory: 0.0,
				CPU:    0.0,
			})
			continue
		}

		proc.Memory = memoryInfo.RSS

		cpuInfo, err := p.CPUPercent()
		if err != nil {
			proc.CPU = 0.0
		}

		proc.CPU = cpuInfo

		processesInfo = append(processesInfo, proc)
	}

	// Getting only "n" number of processes
	if len(processesInfo) > n {
		processesInfo = processesInfo[:n]
	}

	// Sorting processes by CPU usage
	sort.Slice(processesInfo, func(i, j int) bool {
		return processesInfo[i].CPU > processesInfo[j].CPU
	})

	return processesInfo, nil
}
