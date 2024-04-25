package services

import (
	"context"
	"time"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/internals/storage"
	"github.com/abielalejandro/shortener-service/pkg/logger"
	"github.com/abielalejandro/shortener-service/pkg/utils"
)

type Service interface {
	GenerateShort(url string) (string, error)
	SearchUrlByShort(url string) (string, error)
}

type ShortenerService struct {
	config       *config.Config
	log          *logger.Logger
	storage      storage.Storage
	cachestorage storage.CacheStorage
}

func NewShortenerService(config *config.Config,
	storage storage.Storage,
	cachestorage storage.CacheStorage) *ShortenerService {

	return &ShortenerService{
		config:       config,
		log:          logger.New(config.Log.Level),
		storage:      storage,
		cachestorage: cachestorage,
	}
}

func (svc *ShortenerService) GenerateShort(longUrl string) (string, error) {
	ctx := context.Background()
	exists, err := svc.cachestorage.ExistsByFilter(ctx, svc.config.CacheStorage.FilterName, longUrl)
	if err != nil {
		return "", err
	}
	short := utils.ToBase62(longUrl)
	if !exists {
		url := &storage.Url{Long: longUrl, Short: short, CreatedAt: time.Now(), LastVisited: time.Now(), ExpiresAt: time.Now()}
		_, err = svc.storage.Create(ctx, url)

		if err != nil {
			return "", err
		}
		svc.cachestorage.AddFilter(ctx, svc.config.CacheStorage.FilterName, longUrl)
		svc.cachestorage.Add(ctx, short, longUrl, svc.config.CacheStorage.ExpireTimeInMinutes)
	}

	return short, nil
}

func (svc *ShortenerService) SearchUrlByShort(short string) (string, error) {
	ctx := context.Background()
	long, err := svc.cachestorage.Get(ctx, short)
	if err == nil {
		return long, nil
	}
	url, err := svc.storage.GetUrlByShort(ctx, short)

	if err != nil {
		return "", err
	}
	url.LastVisited = time.Now()
	svc.storage.Update(ctx, url)
	svc.cachestorage.Add(ctx, short, url.Long, svc.config.CacheStorage.ExpireTimeInMinutes)
	return url.Long, nil
}
