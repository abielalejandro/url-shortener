package app

import (
	"github.com/abielalejandro/tgs-service/config"
	"github.com/abielalejandro/tgs-service/internals/services"
	"github.com/abielalejandro/tgs-service/internals/storage"
	"github.com/abielalejandro/tgs-service/pkg/logger"
)

type App struct {
	storage.Storage
	*config.Config
	*services.TgsService
}

func NewApp(config *config.Config) *App {
	return &App{
		Config:     config,
		Storage:    storage.NewStorage(config),
		TgsService: services.NewTgsService(config),
	}
}

func (app *App) Run() {
	l := logger.New(app.Config.Log.Level)
	l.Info("App Running")

	app.TgsService.GenerateRange(app.Storage)
	for i := 0; i < 1000; i++ {
		next, err := app.TgsService.GenerateToken()
		if err != nil {
			app.TgsService.GenerateRange(app.Storage)
		}
		l.Info(next)
	}

}
