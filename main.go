package main

import (
	"fmt"
	"if-win-dex-agent/internal/pdh"
)

type thermalZoneData struct {
	Name                     string
	HighPrecisionTemperature float64 `perfdata:"High Precision Temperature"`
	PercentPassiveLimit      float64 `perfdata:"% Passive Limit"`
	ThrottleReasons          float64 `perfdata:"Throttle Reasons"`
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

func main() {
	var thermalZoneDataResults []thermalZoneData
	var tcpDataResults []tcpData

	thermalZoneDataCollector, err := pdh.NewCollector[thermalZoneData]("Thermal Zone Information", pdh.InstancesAll)
	tcpDataCollector, err := pdh.NewCollector[tcpData]("TCPv4", pdh.InstancesAll)
	if err != nil {
		println(err.Error())
	}
	err = thermalZoneDataCollector.Collect(&thermalZoneDataResults)
	if err != nil {
		println(err.Error())
	}
	err = tcpDataCollector.Collect(&tcpDataResults)
	fmt.Println(thermalZoneDataResults[0])
	println(tcpDataResults[0].ConnectionFailures)

}
