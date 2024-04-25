package storage

import (
	"context"

	"github.com/abielalejandro/shortener-service/config"
)

type NotFoundError struct{}

func (err *NotFoundError) Error() string {
	return "not found"
}

type Storage interface {
	ExistsByShort(ctx context.Context, key string) (bool, error)
	Create(ctx context.Context, url *Url) (bool, error)
	Update(ctx context.Context, url *Url) (bool, error)
	GetUrlByShort(ctx context.Context, shortUrl string) (*Url, error)
	GetUrlByLong(ctx context.Context, longUrl string) (*Url, error)
}

func NewStorage(config *config.Config) Storage {
	switch config.Storage.Type {
	case "generic":
		return NewGenericStorage(config)
	case "cassandra":
		return NewCassandraStorage(config)
	default:
		return NewGenericStorage(config)
	}
}

type CacheStorage interface {
	Add(ctx context.Context, key string, val string, minutesToExpire int) error
	AddFilter(ctx context.Context, key string, val string) error
	Exists(ctx context.Context, key string) (bool, error)
	Get(ctx context.Context, key string) (string, error)
	ExistsByFilter(ctx context.Context, key string, val string) (bool, error)
	IsRateLimiterValid(ctx context.Context, key string, limit int, minutes int) (bool, error)
}

func NewCacheStorage(config *config.Config) CacheStorage {

	switch config.CacheStorage.Type {
	case "generic":
		return NewGenericCacheStorage(config)
	case "redis":
		return NewRedisStorage(config)
	default:
		return NewGenericCacheStorage(config)
	}
}
