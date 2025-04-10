package main

import (
	"github.com/cmpxNot29a/shurs/internal/app"
	"github.com/cmpxNot29a/shurs/internal/config"
	"log"
)

func main() {

	conf := config.LoadConfig()

	log.Printf("Configuration loaded: ServerAddress=%s, BaseURL=%s", conf.ServerAddress, conf.BaseURL)

	if err := app.App(conf); err != nil {
		log.Fatalf("FATAL: Application run failed: %v", err)
	}
}
