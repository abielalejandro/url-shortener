package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App     `yaml:"app"`
		Api     `yaml:"api"`
		HTTP    `yaml:"http"`
		RPC     `yaml:"rcp"`
		Log     `yaml:"logger"`
		Storage `yaml:"storage"`
		Token   `yaml:"token"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	RPC struct {
		Port string `env-required:"true" yaml:"port" env:"RPC_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	Storage struct {
		Type  string `env-required:"true" yaml:"type"  env:"STORAGE_TYPE" env-default: "generic"`
		Redis Redis  `yaml:"redis"`
	}

	Api struct {
		Type string `env-required:"true" yaml:"type"  env:"API_TYPE" env-default: "http"`
	}

	// Redis -.
	Redis struct {
		Addr         string `yaml:"addr" env:"REDIS_HOST" env-default:"localhost:6379"`
		Password     string `env:"REDIS_PWD" env-default:""`
		Db           int    `env:"REDIS_DB" env-default:"0"`
		SequenceName string `env:"REDIS_SEQUENCE_NAME" env-default:"tgs"`
	}
	Token struct {
		Range int `yaml:"range" env: "TOKEN_RANGE" env-default:"6"`
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
