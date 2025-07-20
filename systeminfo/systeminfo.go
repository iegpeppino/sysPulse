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
// Will be used in future implementations
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
		return 0.0, err
	}

	return cpuPercentage[0], nil
}

// Return Cpu Loads
func GetCPULoad() ([]cpu.TimesStat, error) {

	cpuLoad, err := cpu.Times(false)
	if err != nil {
		return []cpu.TimesStat{}, err
	}

	currLoad := cpuLoad[0]

	// Calculate total CPU load
	totalLoad := currLoad.Guest + currLoad.Idle + currLoad.Iowait + currLoad.Irq +
		currLoad.Nice + currLoad.Softirq + currLoad.Steal + currLoad.System + currLoad.User

	// Convert loads to percentual values using totalLoad
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
		return &mem.VirtualMemoryStat{}, err
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

// Disk partition stats struct
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
		// Use mountpoint to get disks
		// Virtual memory returns filepaths
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
		return disks, errors.New("disks couldn't be found")
	}

	// Sort Disk by total capacity
	sort.Slice(disks, func(i, j int) bool {
		return disks[i].Total > disks[j].Total
	})

	return disks, nil
}

// Running process stats struct
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
		return []ProcessInfo{}, err
	}

	var processesInfo []ProcessInfo
	var procErr error

	for _, p := range processes {
		proc := ProcessInfo{}
		proc.PID = p.Pid

		proc.Name, err = p.Name()
		if err != nil {
			procErr = errors.Join(procErr, err)
			proc.Name = "N/A"
		}

		proc.Status, err = p.Status()
		if err != nil {
			procErr = errors.Join(procErr, err)
			proc.Status = []string{"Unknown"}
		}

		started, err := p.CreateTime()
		if err != nil {
			procErr = errors.Join(procErr, err)
			proc.Runtime = "N/A"
		}

		// Divide by 1000 since CreateTime() returns uint time in milliseconds
		runtime := time.Since(time.Unix(started/1000, 0)).Truncate(time.Second)
		proc.Runtime = runtime.String()

		cpuInfo, err := p.CPUPercent()
		if err != nil {
			procErr = errors.Join(procErr, err)
			proc.CPU = 0.0
		}

		proc.CPU = cpuInfo

		memoryInfo, err := p.MemoryInfo()
		// If for loop is not broken after a memoryInfo error
		// a runtime error occurs
		if err != nil {
			procErr = errors.Join(procErr, err)
			processesInfo = append(processesInfo, ProcessInfo{
				PID:     proc.PID,
				Name:    proc.Name,
				Status:  proc.Status,
				Runtime: proc.Runtime,
				Memory:  0.0,
				CPU:     0.0,
			})
			continue
		}

		proc.Memory = memoryInfo.RSS

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

	return processesInfo, procErr
}
