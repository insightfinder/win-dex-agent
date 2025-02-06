package collector

import (
	"github.com/shirou/gopsutil/disk"
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

	"golang.org/x/sys/windows"
)

type GeneralCollector struct {
}

func CreateGeneralCollector() *GeneralCollector {
	return &GeneralCollector{}
}

func (collector *GeneralCollector) GetMemoryMetrics() *map[string]float64 {
	result := make(map[string]float64)
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("Error fetching virtual memory info: %v", err)
		return &result
	}

	result["Memory Available MB"] = float64(vmStat.Available) / 1024 / 1024
	result["Memory Used MB"] = float64(vmStat.Used) / 1024 / 1024
	result["Memory Usage %"] = vmStat.UsedPercent

	return &result
}

func (collector *GeneralCollector) GetCPUMetrics() *map[string]float64 {
	result := make(map[string]float64)

	cpuPercentData, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatalf("Error fetching CPU percent data: %v", err)
		return &result
	}
	result["CPU Usage %"] = cpuPercentData[0]

	return &result
}

func (collector *GeneralCollector) GetDiskMetrics() *map[string]float64 {
	result := make(map[string]float64)
	// Call GetLogicalDriveStrings to get the list of drives
	buffer := make([]uint16, 256)
	ret, err := windows.GetLogicalDriveStrings(uint32(len(buffer)), &buffer[0])
	if err != nil {
		return &result
	}

	// Convert the buffer to a list of drive strings
	drives := make([]string, 0)
	for i := 0; i < int(ret); i += 4 {
		drive := windows.UTF16ToString(buffer[i/2 : i/2+2])
		if drive != "" {
			drives = append(drives, drive)
		}
	}

	// Get disk usage for each drive
	for _, drive := range drives {
		usage, err := disk.Usage(drive)
		if err != nil {
			log.Fatalf("Error fetching disk usage for %s: %v", drive, err)
			continue
		}
		//result["Disk "+drive+" Total"] = float64(usage.Total) / 1024 / 1024
		//result["Disk "+drive+" Used"] = float64(usage.Used) / 1024 / 1024
		//result["Disk "+drive+" Free"] = float64(usage.Free) / 1024 / 1024
		result[drive+"\\ "+"Usage %"] = usage.UsedPercent
	}

	return &result
}
