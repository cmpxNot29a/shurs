package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	ServerAddress string // Адрес запуска HTTP-сервера
	BaseURL       string // Базовый адрес для сокращенных URL
}

func LoadConfig() *Config {

	cfg := &Config{}

	const defaultServerAddress = ":8080"
	const defaultBaseURL = "http://localhost:8080"

	flag.StringVar(&cfg.ServerAddress, "a", defaultServerAddress, "HTTP server start address")
	flag.StringVar(&cfg.BaseURL, "b", defaultBaseURL, "Base address for resulting short URLs")

	flag.Parse()

	if envVar := os.Getenv("SERVER_ADDRESS"); envVar != "" {
		cfg.ServerAddress = envVar
	}

	if envVar := os.Getenv("BASE_URL"); envVar != "" {
		cfg.BaseURL = envVar
	}

	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")

	return cfg
}
