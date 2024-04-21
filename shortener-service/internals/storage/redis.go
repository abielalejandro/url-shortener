package storage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	config *config.Config
	client *redis.Client
	log    *logger.Logger
}

func (storage *RedisStorage) Add(key string, val string) error {
	ctx := context.Background()

	next, err := storage.client.Set(ctx, key, val, redis.KeepTTL).Result()
	if err != nil {
		storage.log.Error("Add")
		storage.log.Error(err)
		return err
	}

	storage.log.Info(fmt.Sprintf("GetNext %v", next))

	return nil
}

func (storage *RedisStorage) AddFilter(key string, val string) error {
	ctx := context.Background()

	next, err := storage.client.BFAdd(ctx, key, val).Result()
	if err != nil {
		storage.log.Error("BFAdd")
		storage.log.Error(err)
		return err
	}

	storage.log.Info(fmt.Sprintf("GetNext %v", next))

	return nil
}

func (storage *RedisStorage) Exists(key string) (bool, error) {
	ctx := context.Background()

	next, err := storage.client.Exists(ctx, key).Result()
	if err != nil {
		storage.log.Error("Exists:Exists")
		storage.log.Error(err)
		return false, err
	}

	return next > 0, nil
}

func (storage *RedisStorage) ExistsByFilter(key string) (bool, error) {
	ctx := context.Background()

	var val string
	next, err := storage.client.BFExists(ctx, key, val).Result()
	if err != nil {
		storage.log.Error("ExistsByFilter:BFExists")
		storage.log.Error(err)
		return false, err
	}

	return next, nil
}

func (storage *RedisStorage) saveLimiter(ctx context.Context, key string, duration time.Duration) error {
	_, err := storage.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {

		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, duration)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (storage *RedisStorage) IsRateLimiterValid(key string, limit int, minutes int) (bool, error) {
	ctx := context.Background()

	var err error
	var sequence string
	var val int

	duration, _ := time.ParseDuration(fmt.Sprintf("%vm", minutes))
	condition := fmt.Sprintf("%v:%v", key, minutes)
	err = storage.saveLimiter(ctx, condition, duration)
	if err != nil {
		storage.log.Error("IsRateLimiterValid:saveLimiter")
		storage.log.Error(err)
		return false, nil
	}

	sequence, err = storage.client.Get(ctx, condition).Result()
	if err != nil {
		storage.log.Error("IsRateLimiterValid:Get")
		storage.log.Error(err)
		return false, err
	}

	val, err = strconv.Atoi(sequence)
	if err != nil {
		return false, nil
	}

	return val <= limit, nil
}

func NewRedisStorage(config *config.Config) CacheStorage {

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
