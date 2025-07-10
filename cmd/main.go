package main

import (
	"fmt"
	"github/iegpeppino/syspulse/systeminfo"
)

func main() {

	processes, _ := systeminfo.GetProcessInfo() // Different Processes running

	memory, _ := systeminfo.GetMEMLoad() // Memory Load statistics

	cpuInfo, _ := systeminfo.GetCPUinfo() // CPU information

	cpuPercent, _ := systeminfo.GetCPUPercent() // Total CPU Usage

	cpuLoads, _ := systeminfo.GetCPULoad() // Different CPU usage loads

	disk, _ := systeminfo.GetDISKUse() // Disk partitions Statistics

	fmt.Printf("Number of processes running: %d\n", len(processes))
	fmt.Printf("Available memory: %d\n", memory.Available)
	fmt.Printf("Cpu model: %s\n", cpuInfo.ModelName)
	fmt.Printf("CPU usage percentage %f\n", cpuPercent)
	fmt.Printf("System CPU usage: %f\n", cpuLoads.System)
	fmt.Printf("Disk mountpoint: %s\n", disk[0].Partition.Mountpoint)
	fmt.Printf("Disk total: %d\n", disk[0].Total)

}
