package app

import (
	"github.com/abielalejandro/tgs-service/config"
	"github.com/abielalejandro/tgs-service/internals/storage"
	"github.com/abielalejandro/tgs-service/pkg/logger"
)

type App struct {
	storage.Storage
	*config.Config
}

func NewApp(config *config.Config) *App {
	return &App{
		Config:  config,
		Storage: storage.NewStorage(config),
	}
}
func (app *App) Run() {
	l := logger.New(app.Config.Log.Level)

	l.Info("App Running")
}
