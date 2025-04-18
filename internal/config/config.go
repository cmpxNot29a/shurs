package config

import (
	"flag"
	"os"
	"strings"
)

const (
	DefaultServerAddress = ":8080"
	DefaultBaseURL       = "http://localhost:8080"
	DefaultIDLength      = 8
	DefaultAttempts      = 10
)

type Config struct {
	ServerAddress string // Адрес запуска HTTP-сервера
	BaseURL       string // Базовый адрес для сокращенных URL
	IDLength      int
	Attempts      int
}

func LoadConfig() *Config {

	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "a", DefaultServerAddress, "HTTP server start address")
	flag.StringVar(&cfg.BaseURL, "b", DefaultBaseURL, "Base address for resulting short URLs")

	flag.Parse()

	cfg.IDLength = DefaultIDLength
	cfg.Attempts = DefaultAttempts

	if envVar := os.Getenv("SERVER_ADDRESS"); envVar != "" {
		cfg.ServerAddress = envVar
	}

	if envVar := os.Getenv("BASE_URL"); envVar != "" {
		cfg.BaseURL = envVar
	}

	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")

	return cfg
}
