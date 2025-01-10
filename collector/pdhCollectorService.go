package collector

import (
	"if-win-dex-agent/internal/pdh"
	"log/slog"
)

type PdhCollectorService struct {
	CpuData         []cpuData
	MemoryData      []memoryData
	DiskData        []diskData
	ThermalZoneData []thermalZoneData
	TcpData         []tcpData
	NetworkData     []networkData
}

func NewPdhCollectorService() *PdhCollectorService {
	return &PdhCollectorService{}
}

func (p *PdhCollectorService) Collect() {

	var err error
	cpuDataCollector, _ := pdh.NewCollector[cpuData]("Processor Information", pdh.InstancesAll)
	memoryDataCollector, _ := pdh.NewCollector[memoryData]("Memory", pdh.InstancesAll)
	physicalDiskDataCollector, _ := pdh.NewCollector[diskData]("PhysicalDisk", pdh.InstancesAll)
	thermalZoneDataCollector, _ := pdh.NewCollector[thermalZoneData]("Thermal Zone Information", pdh.InstancesAll)
	tcpDataCollector, _ := pdh.NewCollector[tcpData]("TCPv4", pdh.InstancesAll)
	networkDataCollector, _ := pdh.NewCollector[networkData]("Network Interface", pdh.InstancesAll)

	// Collect data
	err = cpuDataCollector.Collect(&p.CpuData)
	if err != nil {
		slog.Error(err.Error())
		p.CpuData = nil
	}

	err = memoryDataCollector.Collect(&p.MemoryData)
	if err != nil {
		slog.Error(err.Error())
		p.MemoryData = nil
	}

	err = physicalDiskDataCollector.Collect(&p.DiskData)
	if err != nil {
		slog.Error(err.Error())
		p.DiskData = nil
	}

	err = thermalZoneDataCollector.Collect(&p.ThermalZoneData)
	if err != nil {
		slog.Error(err.Error())
		p.ThermalZoneData = nil
	}

	err = tcpDataCollector.Collect(&p.TcpData)
	if err != nil {
		slog.Error(err.Error())
		p.TcpData = nil
	}

	err = networkDataCollector.Collect(&p.NetworkData)
	if err != nil {
		slog.Error(err.Error())
		p.NetworkData = nil
	}

}
