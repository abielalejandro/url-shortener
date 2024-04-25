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

func (storage *GenericCacheStorage) Add(ctx context.Context, key string, val string, minutesToExpire int) error {
	storage.cache[key] = val
	return nil
}

func (storage *GenericCacheStorage) AddFilter(ctx context.Context, key string, val string) error {
	storage.filterCache[key] = val
	return nil
}

func (storage *GenericCacheStorage) Exists(ctx context.Context, key string) (bool, error) {
	_, found := storage.cache[key]
	return found, nil
}

func (storage *GenericCacheStorage) ExistsByFilter(ctx context.Context, key string, val string) (bool, error) {
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

func (storage *GenericCacheStorage) IsRateLimiterValid(ctx context.Context, key string, limit int, minutes int) (bool, error) {
	duration, _ := time.ParseDuration(fmt.Sprintf("%vm", minutes))
	condition := fmt.Sprintf("%v:%v", key, minutes)

	storage.saveLimiter(ctx, condition, duration)
	counter := storage.limiter[condition]

	return counter <= limit, nil
}

func (storage *GenericCacheStorage) Get(ctx context.Context, key string) (string, error) {
	found, _ := storage.Exists(ctx, key)
	if !found {
		return "", nil
	}
	return storage.cache[key], nil
}

func NewGenericCacheStorage(config *config.Config) CacheStorage {
	return &GenericCacheStorage{
		Config:      config,
		cache:       make(map[string]string),
		filterCache: make(map[string]string),
		limiter:     make(map[string]int),
	}
}
