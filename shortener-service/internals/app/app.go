package app

import (
	"fmt"

	"github.com/abielalejandro/shortener-service/api"
	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/internals/services"
	"github.com/abielalejandro/shortener-service/internals/storage"
	"github.com/abielalejandro/shortener-service/pkg/logger"
)

type App struct {
	storage.Storage
	*config.Config
	*services.TgsService
	api.Api
}

func NewApp(config *config.Config) *App {
	storage := storage.NewStorage(config)
	svc := services.NewTgsService(config, storage)
	return &App{
		Config:     config,
		Storage:    storage,
		TgsService: svc,
		Api:        api.NewApi(config, svc),
	}
}

func (app *App) Run() {
	l := logger.New(app.Config.Log.Level)
	l.Info(fmt.Sprintf("App Running WITH %s", app.Config.Api.Type))
	l.Info(fmt.Sprintf("Config %v", app.Config))
	app.TgsService.GenerateRange()
	app.Api.Run()
}
