package storage

import "github.com/abielalejandro/tgs-service/config"

type RedisStorage struct {
	Config *config.Config
}

func (storage *RedisStorage) Connect() error {
	return nil
}

func (storage *RedisStorage) GetNext(sequenceName string) (int, error) {
	return 1, nil
}

func NewRedisStorage(config *config.Config) Storage {
	return &RedisStorage{
		Config: config,
	}
}
