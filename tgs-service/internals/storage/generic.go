package storage

import (
	"math/rand"
	"time"

	"github.com/abielalejandro/tgs-service/config"
)

type GenericStorage struct {
	Config *config.Config
}

func (storage *GenericStorage) GetNext(sequenceName string) (int, error) {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(9) + 1, nil
}

func NewGenericStorage(config *config.Config) Storage {
	return &GenericStorage{
		Config: config,
	}
}
