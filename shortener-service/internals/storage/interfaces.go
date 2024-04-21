package storage

import "github.com/abielalejandro/shortener-service/config"

type Storage interface {
	ExistsByShortUrl(key string) (bool, error)
	Create(url string, shortUrl string) (bool, error)
	GetUrlByShortUrl(shortUrl string) (string, error)
}

func NewStorage(config *config.Config) Storage {
	return nil
}

type CacheStorage interface {
	Add(key string, val string, minutesToExpire int) error
	AddFilter(key string, val string) error
	Exists(key string) (bool, error)
	ExistsByFilter(key string) (bool, error)
	IsRateLimiterValid(key string, limit int, minutes int) (bool, error)
}

func NewCacheStorage(config *config.Config) CacheStorage {

	switch config.Storage.Type {
	case "generic":
		return NewGenericCacheStorage(config)
	case "redis":
		return NewRedisStorage(config)
	default:
		return NewGenericCacheStorage(config)
	}
}
