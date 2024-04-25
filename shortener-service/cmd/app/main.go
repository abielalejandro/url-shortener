package main

import (
	"log"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/internals/app"
)

func main() {
	log.Println("Starting app")

	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app := app.NewApp(cfg)
	app.Run()
}
