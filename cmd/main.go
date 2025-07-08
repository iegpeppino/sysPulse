package main

import (
	"fmt"
	"github/iegpeppino/syspulse/systeminfo"
)

func main() {

	processes, _ := systeminfo.GetProcessInfo()

	memory, _ := systeminfo.GetMEMLoad()

	cpuInfo, _ := systeminfo.GetCPUinfo()

	cpuLoads, _ := systeminfo.GetCPULoad()

	disk, _ := systeminfo.GetDISKUse()

	fmt.Printf("Number of processes running: %d\n", len(processes))
	fmt.Printf("Available memory: %d\n", memory.Available)
	fmt.Printf("Cpu model: %s\n", cpuInfo.ModelName)
	fmt.Printf("System CPU usage: %f\n", cpuLoads.System)
	fmt.Printf("Disk mountpoint: %s\n", disk[0].Partition.Mountpoint)
	fmt.Printf("Disk total: %d\n", disk[0].Total)

}
