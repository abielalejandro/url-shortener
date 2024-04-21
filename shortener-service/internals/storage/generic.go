package storage

import (
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
