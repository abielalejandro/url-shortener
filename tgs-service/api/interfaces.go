package api

import (
	"github.com/abielalejandro/tgs-service/config"
	"github.com/abielalejandro/tgs-service/internals/services"
)

type Api interface {
	Run()
}

func NewApi(config *config.Config, svc *services.TgsService) Api {

	switch config.Api.Type {
	case "http":
		return NewHttpApi(config, svc)
	case "grpc":
		return NewRpcApi(config, svc)
	default:
		return NewHttpApi(config, svc)
	}
}
