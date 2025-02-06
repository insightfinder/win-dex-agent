package cache

type Metric struct {
	Metric string `gorm:"primaryKey"`
	Value  float64
}
