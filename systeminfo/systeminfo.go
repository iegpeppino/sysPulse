package systeminfo

import (
	"errors"
	"fmt"

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

// Return Cpu Loads
func GetCPULoad() (cpu.TimesStat, error) {

	cpuLoad, err := cpu.Times(false)
	if err != nil {
		return cpu.TimesStat{}, fmt.Errorf("unable to get CPU times: %w", err)
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

	return currLoad, nil

}

func GetMEMLoad() (mem.VirtualMemoryStat, error) {

	v, err := mem.VirtualMemory()
	if err != nil {
		return mem.VirtualMemoryStat{}, fmt.Errorf("unable to get memory state: %w", err)
	}

	return mem.VirtualMemoryStat{
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
	disks := make([]DiskInfo, len(partitions))
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

	return disks, nil
}

type ProcessInfo struct {
	PID     int32
	Name    string
	CPU     float64
	Memory  float64
	Runtime string
	Status  []string
}

func GetProcessInfo() ([]ProcessInfo, error) {

	processes, err := process.Processes()
	if err != nil {
		return []ProcessInfo{}, fmt.Errorf("unable to read running processes: %w", err)
	}

	procesesInfo := make([]ProcessInfo, len(processes))

	for _, p := range processes {
		proc := ProcessInfo{}
		proc.PID = p.Pid

		proc.Name, err = p.Name()
		if err != nil {
			proc.Name = "N/A"
		}

		memoryInfo, err := p.MemoryInfo()
		if err != nil {
			proc.Memory = 0.0
		}
		proc.Memory = float64(memoryInfo.RSS)

		cpuInfo, err := p.CPUPercent()
		if err != nil {
			proc.CPU = 0.0
		}
		proc.CPU = cpuInfo

		proc.Status, err = p.Status()
		if err != nil {
			proc.Status = []string{"unknown"}
		}

		procesesInfo = append(procesesInfo, proc)
	}

	return procesesInfo, nil
}
