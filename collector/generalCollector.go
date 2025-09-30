package collector

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

type GeneralCollector struct {
}

func CreateGeneralCollector() *GeneralCollector {
	return &GeneralCollector{}
}

func (collector *GeneralCollector) GetMemoryMetrics() *map[string]map[string]float64 {
	result := make(map[string]map[string]float64)
	result[""] = make(map[string]float64)
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("Error fetching virtual memory info: %v", err)
		return &result
	}

	result[""]["Memory Available MB"] = float64(vmStat.Available) / 1024 / 1024
	result[""]["Memory Used MB"] = float64(vmStat.Used) / 1024 / 1024
	result[""]["Memory Usage %"] = vmStat.UsedPercent

	return &result
}

func (collector *GeneralCollector) GetCPUMetrics() *map[string]map[string]float64 {
	result := make(map[string]map[string]float64)
	result[""] = make(map[string]float64)

	cpuPercentData, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatalf("Error fetching CPU percent data: %v", err)
		return &result
	}
	result[""]["CPU Usage %"] = cpuPercentData[0]

	return &result
}

func (collector *GeneralCollector) GetDiskMetrics() *map[string]map[string]float64 {
	samplingTime := 2 * time.Second
	result := make(map[string]map[string]float64)
	// Initial sample of disk I/O counters
	prevCounters, err := disk.IOCounters()
	if err != nil {
		fmt.Println("Error fetching disk IOCounters:", err)
		return &result
	}
	ticker := time.NewTicker(samplingTime)
	defer ticker.Stop()

	// Collect data for only one tick
	<-ticker.C

	// Get the current counters
	currCounters, err := disk.IOCounters()
	if err != nil {
		fmt.Println("Error fetching disk IOCounters:", err)
		return &result
	}

	// Iterate over each device that is present in both samples
	for device, currStat := range currCounters {
		prevStat, exists := prevCounters[device]
		if !exists {
			continue // device may be new
		}

		// Compute differences over the sampling interval
		deltaReadBytes := currStat.ReadBytes - prevStat.ReadBytes
		deltaWriteBytes := currStat.WriteBytes - prevStat.WriteBytes

		deltaReadCount := currStat.ReadCount - prevStat.ReadCount
		deltaWriteCount := currStat.WriteCount - prevStat.WriteCount

		deltaReadTime := currStat.ReadTime - prevStat.ReadTime
		deltaWriteTime := currStat.WriteTime - prevStat.WriteTime

		// Calculate average latencies (ms per operation) over the interval
		var avgReadLatency float64
		if deltaReadCount > 0 {
			avgReadLatency = float64(deltaReadTime) / float64(deltaReadCount)
		}

		var avgWriteLatency float64
		if deltaWriteCount > 0 {
			avgWriteLatency = float64(deltaWriteTime) / float64(deltaWriteCount)
		}

		// Queue length as reported at sampling time.
		queueLength := currStat.IopsInProgress

		// For a 1-second interval these differences are per-second values.
		readBytesPerSec := deltaReadBytes / uint64(samplingTime.Seconds())
		writeBytesPerSec := deltaWriteBytes / uint64(samplingTime.Seconds())

		result[device] = make(map[string]float64)
		result[device]["Read Bytes/s"] = float64(readBytesPerSec)
		result[device]["Write Bytes/s"] = float64(writeBytesPerSec)
		result[device]["Queue Length"] = float64(queueLength)
		result[device]["Read Latency ms"] = avgReadLatency
		result[device]["Write Latency ms"] = avgWriteLatency
	}

	return &result
}

func (collector *GeneralCollector) GetNetworkMetrics() *map[string]map[string]float64 {
	samplingTime := 2 * time.Second
	result := make(map[string]map[string]float64)

	// Initial sample of network I/O counters
	prevCounters, err := net.IOCounters(true) // true for per-interface stats
	if err != nil {
		fmt.Println("Error fetching network IOCounters:", err)
		return &result
	}

	ticker := time.NewTicker(samplingTime)
	defer ticker.Stop()

	// Collect data for only one tick
	<-ticker.C

	// Get the current counters
	currCounters, err := net.IOCounters(true)
	if err != nil {
		fmt.Println("Error fetching network IOCounters:", err)
		return &result
	}

	// Create a map for easier lookup
	prevCounterMap := make(map[string]net.IOCountersStat)
	for _, counter := range prevCounters {
		prevCounterMap[counter.Name] = counter
	}

	// Iterate over each network interface
	for _, currStat := range currCounters {
		prevStat, exists := prevCounterMap[currStat.Name]
		if !exists {
			continue // interface may be new
		}

		// Skip loopback and inactive interfaces
		if strings.Contains(strings.ToLower(currStat.Name), "loopback") ||
			strings.Contains(strings.ToLower(currStat.Name), "isatap") ||
			strings.Contains(strings.ToLower(currStat.Name), "teredo") {
			continue
		}

		// Compute differences over the sampling interval
		deltaBytesRecv := currStat.BytesRecv - prevStat.BytesRecv
		deltaBytesSent := currStat.BytesSent - prevStat.BytesSent
		deltaPacketsRecv := currStat.PacketsRecv - prevStat.PacketsRecv
		deltaPacketsSent := currStat.PacketsSent - prevStat.PacketsSent
		deltaErrIn := currStat.Errin - prevStat.Errin
		deltaErrOut := currStat.Errout - prevStat.Errout
		deltaDropIn := currStat.Dropin - prevStat.Dropin
		deltaDropOut := currStat.Dropout - prevStat.Dropout

		// Calculate per-second values
		bytesRecvPerSec := float64(deltaBytesRecv) / samplingTime.Seconds()
		bytesSentPerSec := float64(deltaBytesSent) / samplingTime.Seconds()
		packetsRecvPerSec := float64(deltaPacketsRecv) / samplingTime.Seconds()
		packetsSentPerSec := float64(deltaPacketsSent) / samplingTime.Seconds()
		errInPerSec := float64(deltaErrIn) / samplingTime.Seconds()
		errOutPerSec := float64(deltaErrOut) / samplingTime.Seconds()
		dropInPerSec := float64(deltaDropIn) / samplingTime.Seconds()
		dropOutPerSec := float64(deltaDropOut) / samplingTime.Seconds()

		// Store metrics for this interface
		result[currStat.Name] = make(map[string]float64)
		result[currStat.Name]["Network Inbound Bytes/s"] = bytesRecvPerSec
		result[currStat.Name]["Network Outbound Bytes/s"] = bytesSentPerSec
		result[currStat.Name]["Network Inbound MB/s"] = bytesRecvPerSec / 1024 / 1024
		result[currStat.Name]["Network Outbound MB/s"] = bytesSentPerSec / 1024 / 1024
		result[currStat.Name]["Network Inbound Packets/s"] = packetsRecvPerSec
		result[currStat.Name]["Network Outbound Packets/s"] = packetsSentPerSec
		result[currStat.Name]["Network Inbound Errors/s"] = errInPerSec
		result[currStat.Name]["Network Outbound Errors/s"] = errOutPerSec
		result[currStat.Name]["Network Inbound Drops/s"] = dropInPerSec
		result[currStat.Name]["Network Outbound Drops/s"] = dropOutPerSec

		// // Total cumulative counters (useful for monitoring)
		// result[currStat.Name]["Network Total Bytes Received"] = float64(currStat.BytesRecv)
		// result[currStat.Name]["Network Total Bytes Sent"] = float64(currStat.BytesSent)
		// result[currStat.Name]["Network Total Packets Received"] = float64(currStat.PacketsRecv)
		// result[currStat.Name]["Network Total Packets Sent"] = float64(currStat.PacketsSent)
		// result[currStat.Name]["Network Total Errors In"] = float64(currStat.Errin)
		// result[currStat.Name]["Network Total Errors Out"] = float64(currStat.Errout)
		// result[currStat.Name]["Network Total Drops In"] = float64(currStat.Dropin)
		// result[currStat.Name]["Network Total Drops Out"] = float64(currStat.Dropout)
	}

	return &result
}

func (collector *GeneralCollector) GetProcessMetrics() *map[string]map[string]float64 {
	result := make(map[string]map[string]float64)
	processes, err := process.Processes()
	if err != nil {
		log.Printf("Error getting processes: %v", err)
	}

	for _, p := range processes {
		// Get the process name
		name, err := p.Name()
		if err != nil {
			continue
		}
		if strings.TrimSpace(name) == "" {
			continue
		}

		// Get CPU usage
		cpuPercent, err := p.CPUPercent()
		if err != nil {
			continue
		}

		// Get memory usage
		memInfo, err := p.MemoryInfo()
		if err != nil {
			continue
		}
		result[name] = make(map[string]float64)
		result[name]["Process CPU Usage %"] = cpuPercent
		result[name]["Process Memory Used MB"] = float64(memInfo.RSS) / 1024 / 1024
	}

	return &result
}
