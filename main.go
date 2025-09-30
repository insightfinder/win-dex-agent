package main

import (
	"context"
	"if-win-dex-agent/cache"
	"if-win-dex-agent/collector"
	"if-win-dex-agent/insightfinder"
	"if-win-dex-agent/tool"
	"log/slog"
	"time"
)

func main() {
	cacheService, err := cache.CreateCacheService()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if cacheService == nil {
		slog.Error("Failed to create cache service")
		return
	}

	// Init InsightFinder service
	IFClient := insightfinder.CreateInsightFinderClient("https://stg.insightfinder.com", "ashvat", "fa28bf0a1a0f6e77c6eea446383b36d0b5a884d2", "Win-Dex-Agent-Test")

	generalCollectorService := collector.CreateGeneralCollector()
	pdhCollectorService := collector.NewPdhCollectorService()

	for {
		go func() {
			startTime := time.Now()
			slog.Log(context.Background(), slog.LevelInfo, "Start collecting metrics at", "time", startTime)
			pdhCollectorService.Collect()

			// Add metrics from generalCollectorService
			for device, metrics := range *generalCollectorService.GetMemoryMetrics() {
				for metric, value := range metrics {
					cacheService.AddMetricRecord(device, metric, value)
				}
			}
			for device, metrics := range *generalCollectorService.GetCPUMetrics() {
				for metric, value := range metrics {
					cacheService.AddMetricRecord(device, metric, value)
				}
			}
			for device, metrics := range *generalCollectorService.GetProcessMetrics() {
				for metric, value := range metrics {
					cacheService.AddMetricRecord(device, metric, value)
				}
			}
			for device, metrics := range *generalCollectorService.GetNetworkMetrics() {
				for metric, value := range metrics {
					cacheService.AddMetricRecord(device, metric, value)
				}
			}

			// Add metrics from pdhCollectorService
			for device, metrics := range *pdhCollectorService.GetThermalMetrics() {
				for metric, value := range metrics {
					cacheService.AddMetricRecord(device, metric, value)
				}
			}
			for device, metrics := range *pdhCollectorService.GetNetworkMetrics() {
				for metric, value := range metrics {
					cacheService.AddMetricRecord(device, metric, value)
				}
			}

			for device, metrics := range *pdhCollectorService.GetDiskMetrics() {
				for metric, value := range metrics {
					cacheService.AddMetricRecord(device, metric, value)
				}
			}

			idm := tool.BuildIDMFromCache(startTime, "Win-Dex-Agent", cacheService)
			IFClient.SendMetricData(idm)
			cacheService.ClearCache()
			slog.Log(context.Background(), slog.LevelInfo, "End collecting metrics at", "time", time.Now())
		}()
		time.Sleep(5 * time.Minute)
	}

}
