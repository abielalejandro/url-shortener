package api

import (
	"fmt"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/internals/services"
)

type Api interface {
	Run()
}

func NewApi(config *config.Config, svc services.Service, rate *services.RateService) Api {
	fmt.Println(config.Api.Type)
	switch config.Api.Type {
	case "http":
		return NewHttpApi(config, svc, rate)
	case "grpc":
		return NewGrpcApi(config, svc, rate)
	default:
		return NewHttpApi(config, svc, rate)
	}
}
