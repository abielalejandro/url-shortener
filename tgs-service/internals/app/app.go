package app

import (
	"fmt"

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
	api.Api
}

func NewApp(config *config.Config) *App {
	svc := services.NewTgsService(config)
	return &App{
		Config:     config,
		Storage:    storage.NewStorage(config),
		TgsService: svc,
		Api:        api.NewApi(config, svc),
	}
}

func (app *App) Run() {
	l := logger.New(app.Config.Log.Level)
	l.Info(fmt.Sprintf("App Running WITH %s", app.Config.Api.Type))

	app.TgsService.GenerateRange(app.Storage)
	app.Api.Run()
}
