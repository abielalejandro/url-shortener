package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/abielalejandro/shortener-service/config"
)

type GenericCacheStorage struct {
	Config      *config.Config
	cache       map[string]string
	filterCache map[string]string
	limiter     map[string]int
}

func (storage *GenericCacheStorage) Add(key string, val string) error {
	storage.cache[key] = val
	return nil
}

func (storage *GenericCacheStorage) AddFilter(key string, val string) error {
	storage.filterCache[key] = val
	return nil
}

func (storage *GenericCacheStorage) Exists(key string) (bool, error) {
	_, found := storage.cache[key]
	return found, nil
}

func (storage *GenericCacheStorage) ExistsByFilter(key string) (bool, error) {
	_, found := storage.filterCache[key]
	return found, nil
}

func (storage *GenericCacheStorage) saveLimiter(ctx context.Context, key string, duration time.Duration) error {
	val, found := storage.limiter[key]
	if !found {
		val = 0
	}
	val = val + 1
	storage.limiter[key] = val
	return nil
}

func (storage *GenericCacheStorage) IsRateLimiterValid(key string, limit int, minutes int) (bool, error) {
	ctx := context.Background()
	duration, _ := time.ParseDuration(fmt.Sprintf("%vm", minutes))
	condition := fmt.Sprintf("%v:%v", key, minutes)

	storage.saveLimiter(ctx, condition, duration)
	counter := storage.limiter[condition]

	return counter <= limit, nil
}

func NewGenericCacheStorage(config *config.Config) CacheStorage {
	return &GenericCacheStorage{
		Config:      config,
		cache:       make(map[string]string),
		filterCache: make(map[string]string),
		limiter:     make(map[string]int),
	}
}
