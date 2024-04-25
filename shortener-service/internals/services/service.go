package services

import (
	"context"
	"time"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/internals/storage"
	"github.com/abielalejandro/shortener-service/pkg/logger"
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
	tgsService   TgsService
}

func NewShortenerService(
	config *config.Config,
	storage storage.Storage,
	cachestorage storage.CacheStorage,
	tgsService TgsService,
) *ShortenerService {

	return &ShortenerService{
		config:       config,
		log:          logger.New(config.Log.Level),
		storage:      storage,
		cachestorage: cachestorage,
		tgsService:   tgsService,
	}
}

func (svc *ShortenerService) GenerateShort(longUrl string) (string, error) {
	ctx := context.Background()
	var short string
	exists, err := svc.cachestorage.ExistsByFilter(ctx, svc.config.CacheStorage.FilterName, longUrl)
	if err != nil {
		return "", err
	}

	if !exists {
		short, err = svc.tgsService.Next(longUrl)
		if err != nil {
			return "", err
		}

		url := &storage.Url{Long: longUrl, Short: short, CreatedAt: time.Now(), LastVisited: time.Now(), ExpiresAt: time.Now()}
		_, err = svc.storage.Create(ctx, url)

		if err != nil {
			return "", err
		}
		svc.cachestorage.AddFilter(ctx, svc.config.CacheStorage.FilterName, longUrl)
		svc.cachestorage.Add(ctx, short, longUrl, svc.config.CacheStorage.ExpireTimeInMinutes)
	} else {
		url, err := svc.storage.GetUrlByLong(ctx, longUrl)
		if err != nil {
			return "", err
		}
		short = url.Short
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
