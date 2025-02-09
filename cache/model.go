package cache

type Metric struct {
	Device string `gorm:"primaryKey"`
	Metric string `gorm:"primaryKey"`
	Value  float64
}
