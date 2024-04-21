package storage

import (
	"context"
	"fmt"

	"github.com/abielalejandro/tgs-service/config"
	"github.com/abielalejandro/tgs-service/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	config *config.Config
	client *redis.Client
	log    *logger.Logger
}

func (storage *RedisStorage) GetNext(sequenceName string) (int, error) {
	ctx := context.Background()

	next, err := storage.client.Incr(ctx, storage.config.Redis.SequenceName).Result()
	if err != nil {
		storage.log.Error("GetNext:Get")
		storage.log.Error(err)
		return 0, err
	}

	storage.log.Info(fmt.Sprintf("GetNext %v", next))

	return int(next), nil
}

func NewRedisStorage(config *config.Config) Storage {

	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.Db,
	})

	return &RedisStorage{
		config: config,
		client: client,
		log:    logger.New(config.Log.Level),
	}
}
