package storage

import (
	"context"
	"strings"
	"time"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/gocql/gocql"
)

type CassandraStorage struct {
	*gocql.ClusterConfig
}

func NewCassandraStorage(config *config.Config) *CassandraStorage {
	addr := config.Storage.Addr
	hosts := strings.Split(addr, ",")
	clusterConfig := gocql.NewCluster(hosts...)
	clusterConfig.Keyspace = config.Storage.Db

	return &CassandraStorage{
		ClusterConfig: clusterConfig,
	}
}

func (storage *CassandraStorage) GetUrlByLong(ctx context.Context, longUrl string) (*Url, error) {
	return &Url{
		Short:       "12345678",
		Long:        "http://google.com",
		LastVisited: time.Now(),
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now(),
	}, nil
}

func (storage *CassandraStorage) ExistsByShort(ctx context.Context, shortUrl string) (bool, error) {
	return false, nil
}

func (storage *CassandraStorage) Create(ctx context.Context, url *Url) (bool, error) {
	return true, nil
}

func (storage *CassandraStorage) Update(ctx context.Context, url *Url) (bool, error) {
	return true, nil
}

func (storage *CassandraStorage) GetUrlByShort(ctx context.Context, shortUrl string) (*Url, error) {
	return &Url{
		Short:       "12345678",
		Long:        "http://google.com",
		LastVisited: time.Now(),
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now(),
	}, nil
}
