package app

import (
	"github.com/abielalejandro/tgs-service/api"
	"github.com/abielalejandro/tgs-service/config"
	"github.com/abielalejandro/tgs-service/internals/services"
	"github.com/abielalejandro/tgs-service/internals/storage"
	"github.com/abielalejandro/tgs-service/pkg/logger"
)

type App struct {
	storage.Storage
	*config.Config
	*services.TgsService
	*api.HttpApi
}

func NewApp(config *config.Config) *App {
	svc := services.NewTgsService(config)
	return &App{
		Config:     config,
		Storage:    storage.NewStorage(config),
		TgsService: svc,
		HttpApi:    api.NewHttpApi(config, svc),
	}
}

func (app *App) Run() {
	l := logger.New(app.Config.Log.Level)
	l.Info("App Running")

	app.TgsService.GenerateRange(app.Storage)
	app.HttpApi.Run()
}
