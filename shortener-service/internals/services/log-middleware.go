package services

import (
	"fmt"
	"time"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/pkg/logger"
)

type logMiddleware struct {
	next   Service
	config *config.Config
	log    *logger.Logger
}

func (middleware *logMiddleware) GenerateShort(url string) (short string, err error) {
	defer func(start time.Time) {
		middleware.log.Info(fmt.Sprintf("Executing  GenerateShort %s takes %v", short, time.Since(start)))
	}(time.Now())

	return middleware.next.GenerateShort(url)
}

func (middleware *logMiddleware) SearchUrlByShort(short string) (url string, err error) {
	defer func(start time.Time) {
		middleware.log.Info(fmt.Sprintf("Executing  SearchUrlByShort %s takes %v", url, time.Since(start)))
	}(time.Now())

	return middleware.next.SearchUrlByShort(short)

}

func NewLogMiddleware(config *config.Config, next Service) Service {
	return &logMiddleware{
		next:   next,
		config: config,
		log:    logger.New(config.Log.Level),
	}
}
