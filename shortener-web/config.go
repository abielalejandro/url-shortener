package main

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App                  `yaml:"app"`
		HTTP                 `yaml:"http"`
		RestShortenerService `yaml:"shortener_service"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		Domain  string `env-required:"true" yaml:"domain" env:"APP_DOMAIN"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	RestShortenerService struct {
		Url     string `yaml:"addr" env:"SHORTENER_SERVICE_URL" env-default:"http://localhost:8080"`
		Version string `yaml:"version" env:"SHORTENER_SERVICE_VERSION" env-default:"v1"`
		Type    string `env-required:"true" yaml:"type"  env:"SHORTENER_SERVICE_TYPE" env-default: "http"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
