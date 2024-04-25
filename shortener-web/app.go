package main

import (
	"log"
)

type Appl struct {
	*ShortenerService
	*HttpApi
	*Config
}

func NewApp(config *Config) *Appl {
	svc := NewShortenerService(config)
	return &Appl{
		Config:           config,
		ShortenerService: svc,
		HttpApi:          NewHttpApi(config, svc),
	}
}

func (app *Appl) Run() {
	app.HttpApi.Run()
}

func main() {

	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	app := NewApp(cfg)
	app.Run()
}
