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
	storage.CacheStorage
	*config.Config
	*services.ShortenerService
	api.Api
}

func NewApp(config *config.Config) *App {
	db := storage.NewStorage(config)
	cache := storage.NewCacheStorage(config)
	svc := services.NewShortenerService(config, db, cache)
	rate := services.NewRateLimiterService(cache)
	return &App{
		Config:           config,
		Storage:          db,
		CacheStorage:     cache,
		ShortenerService: svc,
		Api:              api.NewApi(config, svc, rate),
	}
}

func (app *App) Run() {
	l := logger.New(app.Config.Log.Level)
	l.Info(fmt.Sprintf("Config %v", app.Config))
	app.Api.Run()
}
