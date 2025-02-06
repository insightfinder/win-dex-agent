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
	IFClient := insightfinder.CreateInsightFinderClient("https://stg.insightfinder.com", "maoyuwang", "", "maoyu-test-win-dex-1")

	generalCollectorService := collector.CreateGeneralCollector()
	pdhCollectorService := collector.NewPdhCollectorService()

	for {
		go func() {
			startTime := time.Now()
			slog.Log(context.Background(), slog.LevelInfo, "Start collecting metrics at", "time", startTime)
			pdhCollectorService.Collect()
			// Add metrics from generalCollectorService
			for _, getMetrics := range []func() *map[string]float64{
				generalCollectorService.GetCPUMetrics,
				generalCollectorService.GetMemoryMetrics,
				generalCollectorService.GetDiskMetrics,
			} {
				for metricName, metricValue := range *getMetrics() {
					cacheService.AddMetricRecord(metricName, metricValue)
				}
			}

			// Add metrics from pdhCollectorService
			for _, getMetrics := range []func() *map[string]float64{
				pdhCollectorService.GetThermalMetrics,
				pdhCollectorService.GetNetworkMetrics,
				pdhCollectorService.GetDiskMetrics,
			} {
				for metricName, metricValue := range *getMetrics() {
					cacheService.AddMetricRecord(metricName, metricValue)
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
