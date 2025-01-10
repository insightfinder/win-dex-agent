package insightfinder

type InstanceDataMap map[string]InstanceData

type InstanceData struct {
	InstanceName       string                    `json:"in" validate:"required"`
	ComponentName      string                    `json:"cn,omitempty"`
	DataInTimestampMap map[int64]DataInTimestamp `json:"dit" validate:"required"`
}
type DataInTimestamp struct {
	TimeStamp        int64             `json:"t" validate:"required"`
	MetricDataPoints []MetricDataPoint `json:"metricDataPointSet" validate:"required"`
}

type MetricDataPoint struct {
	MetricName string  `json:"m" validate:"required"`
	Value      float64 `json:"v" validate:"required"`
}

type MetricDataReceivePayload struct {
	ProjectName      string                  `json:"projectName" validate:"required"`
	UserName         string                  `json:"userName" validate:"required"`
	InstanceDataMap  map[string]InstanceData `json:"idm" validate:"required"`
	SystemName       string                  `json:"systemName,omitempty"`
	MinTimestamp     int64                   `json:"mi,omitempty"`
	MaxTimestamp     int64                   `json:"ma,omitempty"`
	InsightAgentType string                  `json:"iat,omitempty"`
	SamplingInterval string                  `json:"si,omitempty"`
	CloudType        string                  `json:"ct,omitempty"`
}
