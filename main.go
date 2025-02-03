package main

import (
	"if-win-dex-agent/cache"
	"if-win-dex-agent/collector"
	"if-win-dex-agent/insightfinder"
	"if-win-dex-agent/tool"
	"log/slog"
	"strconv"
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

	collectorService := collector.NewPdhCollectorService()
	collectorService.Collect()

	// Add Thermal data
	if collectorService.ThermalZoneData != nil {
		for deviceId, metricData := range collectorService.ThermalZoneData {
			thermalDeviceName := "ThermalZone-" + strconv.FormatInt(int64(deviceId), 10)

			cacheService.AddMetricRecord(thermalDeviceName, "Temperature", (metricData.HighPrecisionTemperature/10.0)-273.15)
			cacheService.AddMetricRecord(thermalDeviceName, "%PassiveLimit", metricData.PercentPassiveLimit)
			cacheService.AddMetricRecord(thermalDeviceName, "ThrottleReason", metricData.ThrottleReasons)

		}
	}

	idm := tool.BuildIDMFromCache(time.Now(), "Win-Dex-Agent", cacheService)
	println(idm)
	IFClient.SendMetricData(idm)
	cacheService.ClearCache()
}
