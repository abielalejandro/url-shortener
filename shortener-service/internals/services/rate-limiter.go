package services

import (
	"context"

	"github.com/abielalejandro/shortener-service/internals/storage"
)

type RateService struct {
	cache storage.CacheStorage
}

func NewRateLimiterService(cache storage.CacheStorage) *RateService {
	return &RateService{cache: cache}
}

func (rateLimiter *RateService) Validate(
	ctx context.Context,
	key string,
	limit int,
	minutes int) (bool, error) {
	return rateLimiter.cache.IsRateLimiterValid(ctx, key, limit, minutes)
}
