package main

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/node_exporter/collector"
)

func main() {
	// Create a new registry
	registry := prometheus.NewRegistry()

	logger := log.NewNopLogger()

	// Create the CPU and memory collectors
	cpuCollector := collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})
	memCollector := collectors.NewGoCollector()

	// Register the collectors with the registry
	registry.MustRegister(cpuCollector)
	registry.MustRegister(memCollector)

	// Enable collectors
	nodeCollector, err := collector.NewNodeCollector(logger)
	if err != nil {
		fmt.Println("Error creating node collector:", err)
		return
	}

	// Register the collector with the registry
	registry.MustRegister(nodeCollector)

	// Gather the metrics
	metricFamilies, err := registry.Gather()
	if err != nil {
		fmt.Println("Error gathering metrics:", err)
		return
	}

	// Print the metrics to stdout
	for _, mf := range metricFamilies {
		fmt.Println(mf)
	}
}
