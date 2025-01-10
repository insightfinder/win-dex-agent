package collector

type cpuData struct {
	Name string

	C1TimeSeconds                   float64 `perfdata:"% C1 Time"`
	C2TimeSeconds                   float64 `perfdata:"% C2 Time"`
	C3TimeSeconds                   float64 `perfdata:"% C3 Time"`
	C1TransitionsTotal              float64 `perfdata:"C1 Transitions/sec"`
	C2TransitionsTotal              float64 `perfdata:"C2 Transitions/sec"`
	C3TransitionsTotal              float64 `perfdata:"C3 Transitions/sec"`
	ClockInterruptsTotal            float64 `perfdata:"Clock Interrupts/sec"`
	DpcQueuedPerSecond              float64 `perfdata:"DPCs Queued/sec"`
	DpcTimeSeconds                  float64 `perfdata:"% DPC Time"`
	IdleBreakEventsTotal            float64 `perfdata:"Idle Break Events/sec"`
	IdleTimeSeconds                 float64 `perfdata:"% Idle Time"`
	InterruptsTotal                 float64 `perfdata:"Interrupts/sec"`
	InterruptTimeSeconds            float64 `perfdata:"% Interrupt Time"`
	ParkingStatus                   float64 `perfdata:"Parking Status"`
	PerformanceLimitPercent         float64 `perfdata:"% Performance Limit"`
	PriorityTimeSeconds             float64 `perfdata:"% Priority Time"`
	PrivilegedTimeSeconds           float64 `perfdata:"% Privileged Time"`
	PrivilegedUtilitySeconds        float64 `perfdata:"% Privileged Utility"`
	ProcessorFrequencyMHz           float64 `perfdata:"Processor Frequency"`
	ProcessorPerformance            float64 `perfdata:"% Processor Performance"`
	ProcessorPerformanceSecondValue float64 `perfdata:"% Processor Performance,secondvalue"`
	ProcessorTimeSeconds            float64 `perfdata:"% Processor Time"`
	ProcessorUtilityRate            float64 `perfdata:"% Processor Utility"`
	ProcessorUtilityRateSecondValue float64 `perfdata:"% Processor Utility,secondvalue"`
	UserTimeSeconds                 float64 `perfdata:"% User Time"`
}

type memoryData struct {
	AvailableBytes                  float64 `perfdata:"Available Bytes"`
	AvailableKBytes                 float64 `perfdata:"Available KBytes"`
	AvailableMBytes                 float64 `perfdata:"Available MBytes"`
	CacheBytes                      float64 `perfdata:"Cache Bytes"`
	CacheBytesPeak                  float64 `perfdata:"Cache Bytes Peak"`
	CacheFaultsPerSec               float64 `perfdata:"Cache Faults/sec"`
	CommitLimit                     float64 `perfdata:"Commit Limit"`
	CommittedBytes                  float64 `perfdata:"Committed Bytes"`
	DemandZeroFaultsPerSec          float64 `perfdata:"Demand Zero Faults/sec"`
	FreeAndZeroPageListBytes        float64 `perfdata:"Free & Zero Page List Bytes"`
	FreeSystemPageTableEntries      float64 `perfdata:"Free System Page Table Entries"`
	ModifiedPageListBytes           float64 `perfdata:"Modified Page List Bytes"`
	PageFaultsPerSec                float64 `perfdata:"Page Faults/sec"`
	PageReadsPerSec                 float64 `perfdata:"Page Reads/sec"`
	PagesInputPerSec                float64 `perfdata:"Pages Input/sec"`
	PagesOutputPerSec               float64 `perfdata:"Pages Output/sec"`
	PagesPerSec                     float64 `perfdata:"Pages/sec"`
	PageWritesPerSec                float64 `perfdata:"Page Writes/sec"`
	PoolNonpagedAllocs              float64 `perfdata:"Pool Nonpaged Allocs"`
	PoolNonpagedBytes               float64 `perfdata:"Pool Nonpaged Bytes"`
	PoolPagedAllocs                 float64 `perfdata:"Pool Paged Allocs"`
	PoolPagedBytes                  float64 `perfdata:"Pool Paged Bytes"`
	PoolPagedResidentBytes          float64 `perfdata:"Pool Paged Resident Bytes"`
	StandbyCacheCoreBytes           float64 `perfdata:"Standby Cache Core Bytes"`
	StandbyCacheNormalPriorityBytes float64 `perfdata:"Standby Cache Normal Priority Bytes"`
	StandbyCacheReserveBytes        float64 `perfdata:"Standby Cache Reserve Bytes"`
	SystemCacheResidentBytes        float64 `perfdata:"System Cache Resident Bytes"`
	SystemCodeResidentBytes         float64 `perfdata:"System Code Resident Bytes"`
	SystemCodeTotalBytes            float64 `perfdata:"System Code Total Bytes"`
	SystemDriverResidentBytes       float64 `perfdata:"System Driver Resident Bytes"`
	SystemDriverTotalBytes          float64 `perfdata:"System Driver Total Bytes"`
	TransitionFaultsPerSec          float64 `perfdata:"Transition Faults/sec"`
	TransitionPagesRePurposedPerSec float64 `perfdata:"Transition Pages RePurposed/sec"`
	WriteCopiesPerSec               float64 `perfdata:"Write Copies/sec"`
}

type diskData struct {
	Name string

	CurrentDiskQueueLength float64 `perfdata:"Current Disk Queue Length"`
	DiskReadBytesPerSec    float64 `perfdata:"Disk Read Bytes/sec"`
	DiskReadsPerSec        float64 `perfdata:"Disk Reads/sec"`
	DiskWriteBytesPerSec   float64 `perfdata:"Disk Write Bytes/sec"`
	DiskWritesPerSec       float64 `perfdata:"Disk Writes/sec"`
	PercentDiskReadTime    float64 `perfdata:"% Disk Read Time"`
	PercentDiskWriteTime   float64 `perfdata:"% Disk Write Time"`
	PercentIdleTime        float64 `perfdata:"% Idle Time"`
	SplitIOPerSec          float64 `perfdata:"Split IO/Sec"`
	AvgDiskSecPerRead      float64 `perfdata:"Avg. Disk sec/Read"`
	AvgDiskSecPerWrite     float64 `perfdata:"Avg. Disk sec/Write"`
	AvgDiskSecPerTransfer  float64 `perfdata:"Avg. Disk sec/Transfer"`
}

type thermalZoneData struct {
	Name                     string
	HighPrecisionTemperature float64 `perfdata:"High Precision Temperature"`
	PercentPassiveLimit      float64 `perfdata:"% Passive Limit"`
	ThrottleReasons          float64 `perfdata:"Throttle Reasons"`
}

type networkData struct {
	Name                     string
	BytesReceivedPerSec      float64 `perfdata:"Bytes Received/sec"`
	BytesSentPerSec          float64 `perfdata:"Bytes Sent/sec"`
	BytesTotalPerSec         float64 `perfdata:"Bytes Total/sec"`
	CurrentBandwidth         float64 `perfdata:"Current Bandwidth"`
	OutputQueueLength        float64 `perfdata:"Output Queue Length"`
	PacketsOutboundDiscarded float64 `perfdata:"Packets Outbound Discarded"`
	PacketsOutboundErrors    float64 `perfdata:"Packets Outbound Errors"`
	PacketsPerSec            float64 `perfdata:"Packets/sec"`
	PacketsReceivedDiscarded float64 `perfdata:"Packets Received Discarded"`
	PacketsReceivedErrors    float64 `perfdata:"Packets Received Errors"`
	PacketsReceivedPerSec    float64 `perfdata:"Packets Received/sec"`
	PacketsReceivedUnknown   float64 `perfdata:"Packets Received Unknown"`
	PacketsSentPerSec        float64 `perfdata:"Packets Sent/sec"`
}

type tcpData struct {
	ConnectionFailures          float64 `perfdata:"Connection Failures"`
	ConnectionsActive           float64 `perfdata:"Connections Active"`
	ConnectionsEstablished      float64 `perfdata:"Connections Established"`
	ConnectionsPassive          float64 `perfdata:"Connections Passive"`
	ConnectionsReset            float64 `perfdata:"Connections Reset"`
	SegmentsPerSec              float64 `perfdata:"Segments/sec"`
	SegmentsReceivedPerSec      float64 `perfdata:"Segments Received/sec"`
	SegmentsRetransmittedPerSec float64 `perfdata:"Segments Retransmitted/sec"`
	SegmentsSentPerSec          float64 `perfdata:"Segments Sent/sec"`
}
