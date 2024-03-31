package storage

import (
	"math/rand"

	"github.com/abielalejandro/tgs-service/config"
)

type GenericStorage struct {
	Config *config.Config
}

func (storage *GenericStorage) GetNext(sequenceName string) (int, error) {
	return rand.Intn(1000000), nil
}

func NewGenericStorage(config *config.Config) Storage {
	return &GenericStorage{
		Config: config,
	}
}
