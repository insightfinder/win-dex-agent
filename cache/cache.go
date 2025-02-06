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

func (cache *CacheService) AddMetricRecord(metric string, value float64) {
	if err := cache.db.Create(&Metric{
		Metric: metric,
		Value:  value,
	}).Error; err != nil {
		slog.Error(err.Error())
	}
}

func (cache *CacheService) GetMetrics() (*[]Metric, error) {
	var metrics []Metric
	err := cache.db.Find(&metrics).Error
	return &metrics, err
}
