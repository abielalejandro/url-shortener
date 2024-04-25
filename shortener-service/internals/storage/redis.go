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

func (cache *RedisStorage) Add(ctx context.Context, key string, val string, minutesToExpire int) error {
	expiration, _ := time.ParseDuration(fmt.Sprintf("%vm", minutesToExpire))
	next, err := cache.client.Set(ctx, key, val, expiration).Result()
	if err != nil {
		cache.log.Error("Add")
		cache.log.Error(err)
		return err
	}

	cache.log.Info(fmt.Sprintf("Add %v", next))

	return nil
}

func (cache *RedisStorage) AddFilter(ctx context.Context, key string, val string) error {
	next, err := cache.client.BFAdd(ctx, key, val).Result()
	if err != nil {
		cache.log.Error("BFAdd")
		cache.log.Error(err)
		return err
	}

	cache.log.Info(fmt.Sprintf("AddFilter %v", next))

	return nil
}

func (cache *RedisStorage) Exists(ctx context.Context, key string) (bool, error) {
	next, err := cache.client.Exists(ctx, key).Result()
	if err != nil {
		cache.log.Error("Exists:Exists")
		cache.log.Error(err)
		return false, err
	}

	return next > 0, nil
}

func (cache *RedisStorage) ExistsByFilter(ctx context.Context, key string, val string) (bool, error) {
	next, err := cache.client.BFExists(ctx, key, val).Result()
	if err != nil {
		cache.log.Error("ExistsByFilter:BFExists")
		cache.log.Error(err)
		return false, err
	}

	return next, nil
}

func (cache *RedisStorage) saveLimiter(ctx context.Context, key string, duration time.Duration) error {
	_, err := cache.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {

		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, duration)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (cache *RedisStorage) IsRateLimiterValid(ctx context.Context, key string, limit int, minutes int) (bool, error) {
	var err error
	var sequence string
	var val int

	duration, _ := time.ParseDuration(fmt.Sprintf("%vm", minutes))
	condition := fmt.Sprintf("%v:%v", key, minutes)
	cache.log.Info(fmt.Sprintf("condition %v", condition))
	cache.log.Info(fmt.Sprintf("duration %v", duration))
	err = cache.saveLimiter(ctx, condition, duration)
	if err != nil {
		cache.log.Error("IsRateLimiterValid:saveLimiter")
		cache.log.Error(err)
		return false, nil
	}

	sequence, err = cache.client.Get(ctx, condition).Result()
	if err != nil {
		cache.log.Error("IsRateLimiterValid:Get")
		cache.log.Error(err)
		return false, err
	}

	cache.log.Info(fmt.Sprintf("sequence %v", sequence))
	val, err = strconv.Atoi(sequence)
	if err != nil {
		return false, nil
	}

	return val <= limit, nil
}

func (cache *RedisStorage) Get(ctx context.Context, key string) (string, error) {
	val, err := cache.client.Get(ctx, key).Result()
	if err != nil {
		cache.log.Error("Get")
		cache.log.Error(err)
		return "", err
	}
	return val, nil
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
