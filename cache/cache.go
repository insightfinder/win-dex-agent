package cache

import (
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
	"log/slog"
)

type CacheService struct {
	db *gorm.DB
}

func CreateCacheService() (*CacheService, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect database")
		return nil, err
	}

	// Auto-migrate the schema for the Metric model
	err = db.AutoMigrate(&Metric{})
	if err != nil {
		return nil, err
	}

	return &CacheService{db: db}, err
}

func (cache *CacheService) ClearCache() {
	err := cache.db.Migrator().DropTable(&Metric{})
	if err != nil {
		slog.Error(err.Error())
	}
	err = cache.db.AutoMigrate(&Metric{})
	if err != nil {
		slog.Error(err.Error())
	}
}

func (cache *CacheService) AddMetricRecord(device string, metric string, value float64) {
	if err := cache.db.Create(&Metric{
		Device: device,
		Metric: metric,
		Value:  value,
	}).Error; err != nil {
		slog.Error(err.Error())
	}
}

func (cache *CacheService) ListDevices() *[]string {
	var devices []string
	if err := cache.db.Model(&Metric{}).
		Distinct("device").
		Pluck("device", &devices).
		Error; err != nil {
		panic(err)
	}
	return &devices
}

func (cache *CacheService) GetMetricsByDevice(device string) *[]Metric {
	var metrics []Metric
	if err := cache.db.Where("device = ?", device).Find(&metrics).Error; err != nil {
		panic(err)
	}
	return &metrics
}
