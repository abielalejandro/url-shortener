package storage

import (
	"github.com/abielalejandro/tgs-service/config"
)

type Storage interface {
	GetNext(sequenceName string) (int, error)
}

func NewStorage(config *config.Config) Storage {

	switch config.Storage.Type {
	case "generic":
		return NewGenericStorage(config)
	case "redis":
		return NewRedisStorage(config)
	default:
		return NewGenericStorage(config)
	}
}
