package api

import (
	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/internals/services"
)

type Api interface {
	Run()
}

func NewApi(config *config.Config, svc services.Service, rate *services.RateService) Api {
	return NewHttpApi(config, svc, rate)
}
