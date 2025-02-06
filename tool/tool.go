package tool

import (
	"if-win-dex-agent/cache"
	"if-win-dex-agent/insightfinder"
	"time"
)

func BuildIDMFromCache(timestamp time.Time, instanceName string, cache *cache.CacheService) *insightfinder.InstanceDataMap {
	instanceDataMap := make(insightfinder.InstanceDataMap)
	metricDataPoints := make([]insightfinder.MetricDataPoint, 0)

	allMetricData, err := cache.GetMetrics()
	if err != nil {
		return nil
	}

	for _, metric := range *allMetricData {
		metricDataPoints = append(metricDataPoints, insightfinder.MetricDataPoint{
			MetricName: metric.Metric,
			Value:      metric.Value,
		})
	}

	dit := make(map[int64]insightfinder.DataInTimestamp)
	dit[timestamp.UnixMilli()] = insightfinder.DataInTimestamp{
		TimeStamp:        timestamp.UnixMilli(),
		MetricDataPoints: metricDataPoints,
	}
	instanceDataMap[instanceName] = insightfinder.InstanceData{
		InstanceName:       instanceName,
		DataInTimestampMap: dit,
	}

	return &instanceDataMap
}
