package api

import (
	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/internals/services"
)

type Api interface {
	Run()
}

func NewApi(config *config.Config, svc services.Service, rate *services.RateService) Api {
	switch config.Api.Type {
	case "http":
		return NewHttpApi(config, svc, rate)
	case "grpc":
		return NewGrpcApi(config, svc)
	default:
		return NewHttpApi(config, svc, rate)
	}
}
