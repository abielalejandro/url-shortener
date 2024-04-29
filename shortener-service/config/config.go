package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App          `yaml:"app"`
		HTTP         `yaml:"http"`
		GRPC         `yaml:"grpc"`
		Log          `yaml:"logger"`
		Storage      `yaml:"storage"`
		CacheStorage `yaml:"cache_storage"`
		RateLimiter  `yaml:"rate_limiter"`
		TgsService   `yaml:"tgs_service"`
		Api          `yaml:"api"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		Domain  string `env-required:"true" yaml:"domain" env:"APP_DOMAIN"`
	}

	GRPC struct {
		Port string `env-required:"true" yaml:"port" env:"GRPC_PORT"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	Storage struct {
		Type     string `yaml:"type" env-required:"true"  env:"STORAGE_TYPE" env-default: "generic"`
		Addr     string `yaml:"addr"  env-required:"true" env:"DB_HOST" env-default:"localhost"`
		Port     string `yaml:"port" env-required:"true"  env:"DB_PORT" env-default:"7000"`
		Password string `yaml:"password" env-required:"true" env:"DB_PWD" env-default:"admin"`
		Db       string `yaml:"db"  env-required:"true" env:"DB_NAME" env-default:"shortener"`
		Username string `yaml:"user" env-required:"true"  env:"DB_USER" env-default:"admin"`
	}

	CacheStorage struct {
		Type                string `env-required:"true" yaml:"type"  env:"CACHE_STORAGE_TYPE" env-default: "generic"`
		ExpireTimeInMinutes int    `env-required:"true" yaml:"expire_time_minutes"  env:"CACHE_STORAGE_EXPIRE_TIME_MINUTES" env-default: "60"`
		FilterName          string `env-required:"true" yaml:"filter_name"  env:"CACHE_STORAGE_FILTER_NAME" env-default: "longurls"`
		Redis               Redis  `yaml:"redis"`
	}

	// Redis -.
	Redis struct {
		Addr     string `yaml:"addr" env:"REDIS_HOST" env-default:"localhost:6379"`
		Password string `env:"REDIS_PWD" env-default:""`
		Db       int    `env:"REDIS_DB" env-default:"0"`
	}

	RateLimiter struct {
		MaxRequests         int `yaml:"max_request" env:"RATE_LIMITER_MAX_REQUEST" env-default:"20"`
		WindowTimeInSeconds int `yaml:"max_request_window_time_seconds" env:"RATE_LIMITER_WINDOW_TIME_SENCONDS" env-default:"60"`
	}

	TgsService struct {
		Url  string `yaml:"addr" env:"TGS_SERVICE_URL" env-default:"http://localhost:8080/api/v1/next"`
		Type string `env-required:"true" yaml:"type"  env:"TGS_SERVICE_TYPE" env-default: "grpc"`
	}

	Api struct {
		Type string `env-required:"true" yaml:"type"  env:"API_TYPE" env-default: "http"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
