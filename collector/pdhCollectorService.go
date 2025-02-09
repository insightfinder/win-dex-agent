package collector

import (
	"if-win-dex-agent/internal/pdh"
	"log/slog"
	"strconv"
	"time"
)

const TicksToSecondScaleFactor = 1 / 1e7

type PdhCollectorService struct {
	memoryData       []memoryData
	diskDataTick1    []diskData
	diskDataTick2    []diskData
	thermalZoneData  []thermalZoneData
	tcpData          []tcpData
	networkDataTick1 []networkData
	networkDataTick2 []networkData
}

func NewPdhCollectorService() *PdhCollectorService {
	return &PdhCollectorService{}
}

func (p *PdhCollectorService) Collect() {

	var err error
	physicalDiskDataCollector, _ := pdh.NewCollector[diskData]("PhysicalDisk", pdh.InstancesAll)
	thermalZoneDataCollector, _ := pdh.NewCollector[thermalZoneData]("Thermal Zone Information", pdh.InstancesAll)
	networkDataCollector, _ := pdh.NewCollector[networkData]("Network Interface", pdh.InstancesAll)

	err = physicalDiskDataCollector.Collect(&p.diskDataTick1)
	if err != nil {
		slog.Error(err.Error())
		p.diskDataTick1 = nil
	}
	time.Sleep(1 * time.Second)
	err = physicalDiskDataCollector.Collect(&p.diskDataTick2)
	if err != nil {
		slog.Error(err.Error())
		p.diskDataTick2 = nil
	}

	err = thermalZoneDataCollector.Collect(&p.thermalZoneData)
	if err != nil {
		slog.Error(err.Error())
		p.thermalZoneData = nil
	}

	err = networkDataCollector.Collect(&p.networkDataTick1)
	if err != nil {
		slog.Error(err.Error())
		p.networkDataTick1 = nil
	}
	time.Sleep(1 * time.Second)
	err = networkDataCollector.Collect(&p.networkDataTick2)
	if err != nil {
		slog.Error(err.Error())
		p.networkDataTick2 = nil
	}

}

func (p *PdhCollectorService) GetDiskMetrics() *map[string]map[string]float64 {
	result := make(map[string]map[string]float64)
	if p.diskDataTick1 == nil || p.diskDataTick2 == nil {
		return &result
	}

	for i, disk1 := range p.diskDataTick1 {
		disk2 := p.diskDataTick2[i]
		diskName := disk1.Name

		// Calculate differences between ticks
		readBytesPerSec := (disk2.DiskReadBytesPerSec - disk1.DiskReadBytesPerSec) / float64(time.Second.Seconds())
		writeBytesPerSec := (disk2.DiskWriteBytesPerSec - disk1.DiskWriteBytesPerSec) / float64(time.Second.Seconds())
		queueLength := disk2.CurrentDiskQueueLength

		// Calculate disk latency
		readLatency := (disk2.AvgDiskSecPerRead - disk1.AvgDiskSecPerRead) * TicksToSecondScaleFactor
		writeLatency := (disk2.AvgDiskSecPerWrite - disk1.AvgDiskSecPerWrite) * TicksToSecondScaleFactor

		result[diskName] = map[string]float64{
			"Read Bytes/s":     readBytesPerSec,
			"Write Bytes/s":    writeBytesPerSec,
			"Queue Length":     queueLength,
			"Read Latency ms":  readLatency,
			"Write Latency ms": writeLatency,
		}
	}
	return &result
}

func (p *PdhCollectorService) GetThermalMetrics() *map[string]map[string]float64 {
	result := make(map[string]map[string]float64)
	if p.thermalZoneData == nil {
		return &result
	}

	for index, thermalZone := range p.thermalZoneData {
		thermalName := "Thermal " + strconv.Itoa(index)
		result[thermalName] = make(map[string]float64)
		result[thermalName]["Temperature"] = (thermalZone.HighPrecisionTemperature / 10.0) - 273.15
		result[thermalName]["Thermal Passive Limit"] = thermalZone.PercentPassiveLimit
		result[thermalName]["Thermal Throttle Reason"] = thermalZone.ThrottleReasons
	}
	return &result
}

func (p *PdhCollectorService) GetNetworkMetrics() *map[string]map[string]float64 {
	result := make(map[string]map[string]float64)
	if p.networkDataTick1 == nil || p.networkDataTick2 == nil {
		return &result
	}

	for i, network1 := range p.networkDataTick1 {
		network2 := p.networkDataTick2[i]
		interfaceName := "Network Interface " + strconv.Itoa(i)

		// Calculate differences between ticks
		receivedBytesPerSec := (network2.BytesReceivedPerSec - network1.BytesReceivedPerSec) / float64(time.Second.Seconds())
		sentBytesPerSec := (network2.BytesSentPerSec - network1.BytesSentPerSec) / float64(time.Second.Seconds())

		result[interfaceName] = map[string]float64{
			"Received MB/s": receivedBytesPerSec / 1024 / 1024,
			"Sent MB/s":     sentBytesPerSec / 1024 / 1024,
		}
	}
	return &result
}
