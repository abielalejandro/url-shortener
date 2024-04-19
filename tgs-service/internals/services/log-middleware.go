package services

import (
	"fmt"
	"time"

	"github.com/abielalejandro/tgs-service/config"
	"github.com/abielalejandro/tgs-service/pkg/logger"
)

type logMiddleware struct {
	next   Service
	config *config.Config
	log    *logger.Logger
}

func (middleware *logMiddleware) GenerateToken() (token string, err error) {
	defer func(start time.Time) {
		middleware.log.Info(fmt.Sprintf("Executing  GenerateToken %s takes %v", token, time.Since(start)))
	}(time.Now())

	return middleware.next.GenerateToken()
}

func NewLogMiddleware(config *config.Config, next Service) Service {
	return &logMiddleware{
		next:   next,
		config: config,
		log:    logger.New(config.Log.Level),
	}
}
