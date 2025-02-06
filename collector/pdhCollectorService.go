package collector

import (
	"if-win-dex-agent/internal/pdh"
	"log/slog"
	"regexp"
	"strconv"
)

const TicksToSecondScaleFactor = 1 / 1e7

type PdhCollectorService struct {
	memoryData      []memoryData
	diskData        []diskData
	thermalZoneData []thermalZoneData
	tcpData         []tcpData
	networkData     []networkData
}

func NewPdhCollectorService() *PdhCollectorService {
	return &PdhCollectorService{}
}

func (p *PdhCollectorService) Collect() {

	var err error
	physicalDiskDataCollector, _ := pdh.NewCollector[diskData]("PhysicalDisk", pdh.InstancesAll)
	thermalZoneDataCollector, _ := pdh.NewCollector[thermalZoneData]("Thermal Zone Information", pdh.InstancesAll)
	networkDataCollector, _ := pdh.NewCollector[networkData]("Network Interface", pdh.InstancesAll)

	err = physicalDiskDataCollector.Collect(&p.diskData)
	if err != nil {
		slog.Error(err.Error())
		p.diskData = nil
	}

	err = thermalZoneDataCollector.Collect(&p.thermalZoneData)
	if err != nil {
		slog.Error(err.Error())
		p.thermalZoneData = nil
	}

	err = networkDataCollector.Collect(&p.networkData)
	if err != nil {
		slog.Error(err.Error())
		p.networkData = nil
	}

}

func (p *PdhCollectorService) GetMemoryMetrics() *map[string]float64 {
	result := make(map[string]float64)
	if p.memoryData == nil {
		return &result
	}
	result["Memory Available"] = p.memoryData[0].AvailableMBytes
	result["Memory Committed"] = p.memoryData[0].CommittedBytes / 1024 / 1024
	result["Memory Used"] = (p.memoryData[0].CommittedBytes - p.memoryData[0].AvailableBytes) / 1024 / 1024
	result["Memory Usage"] = (p.memoryData[0].CommittedBytes - p.memoryData[0].AvailableBytes) / p.memoryData[0].CommittedBytes * 100
	return &result
}

func (p *PdhCollectorService) GetDiskMetrics() *map[string]float64 {
	result := make(map[string]float64)
	var diskName string
	if p.diskData == nil {
		return &result
	}
	for _, disk := range p.diskData {
		re := regexp.MustCompile(`(\d+)\s+([A-Z]):`)
		match := re.FindStringSubmatch(disk.Name)
		if len(match) > 1 {
			diskName = match[2] + ":\\ "
		} else {
			continue
		}
		result[diskName+"Read Bytes/s"] = disk.DiskReadBytesPerSec
		result[diskName+"Write Bytes/s"] = disk.DiskWriteBytesPerSec
		result[diskName+"Queue Length"] = disk.CurrentDiskQueueLength

		// Calculate disk usage percentage
		//totalDiskSpace := disk.PercentFreeSpace * 1024 * 1024
		//freeDiskSpace := disk.FreeSpace * 1024 * 1024
		//result[diskName+"Disk Usage %"] = (totalDiskSpace - freeDiskSpace) / totalDiskSpace * 100

		// Calculate disk latency
		result[diskName+"Read Latency ms"] = disk.AvgDiskSecPerRead * TicksToSecondScaleFactor
		result[diskName+"Write Latency ms"] = disk.AvgDiskSecPerWrite * TicksToSecondScaleFactor
	}
	return &result
}

func (p *PdhCollectorService) GetThermalMetrics() *map[string]float64 {
	result := make(map[string]float64)
	if p.thermalZoneData == nil {
		return &result
	}

	result["Thermal Temperature"] = (p.thermalZoneData[0].HighPrecisionTemperature / 10.0) - 273.15
	result["Thermal Passive Limit"] = p.thermalZoneData[0].PercentPassiveLimit
	result["Thermal Throttle Reason"] = p.thermalZoneData[0].ThrottleReasons
	return &result
}

func (p *PdhCollectorService) GetNetworkMetrics() *map[string]float64 {
	result := make(map[string]float64)
	if p.networkData == nil {
		return &result
	}
	for i, network := range p.networkData {
		interfaceName := "Network Interface " + strconv.Itoa(i)
		result[interfaceName+" Received Bytes/s"] = network.BytesReceivedPerSec
		result[interfaceName+" Sent Bytes/s"] = network.BytesSentPerSec
	}
	return &result
}
