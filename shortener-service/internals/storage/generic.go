package storage

import (
	"context"
	"time"

	"github.com/abielalejandro/shortener-service/config"
)

type GenericStorage struct {
	Config *config.Config
}

func NewGenericStorage(config *config.Config) Storage {
	return &GenericStorage{
		Config: config,
	}
}

func (storage *GenericStorage) GetUrlByLong(ctx context.Context, longUrl string) (*Url, error) {
	return &Url{}, nil
}

func (storage *GenericStorage) ExistsByShort(ctx context.Context, shortUrl string) (bool, error) {
	return false, nil
}

func (storage *GenericStorage) Create(ctx context.Context, url *Url) (bool, error) {
	return true, nil
}

func (storage *GenericStorage) Update(ctx context.Context, url *Url) (bool, error) {
	return true, nil
}

func (storage *GenericStorage) GetUrlByShort(ctx context.Context, shortUrl string) (*Url, error) {
	return &Url{
		Short:       "12345678",
		Long:        "http://google.com",
		LastVisited: time.Now(),
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now(),
	}, nil
}
