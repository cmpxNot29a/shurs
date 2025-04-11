package config

import (
	"flag"
	"strings"
)

type Config struct {
	ServerAddress string // Адрес запуска HTTP-сервера
	BaseURL       string // Базовый адрес для сокращенных URL
}

func LoadConfig() *Config {

	serverAddrPtr := flag.String("a", ":8080", "HTTP server start address")
	baseURLPtr := flag.String("b", "http://localhost:8080", "Base address for resulting short URLs")

	flag.Parse()

	cfg := &Config{
		ServerAddress: *serverAddrPtr,
		BaseURL:       strings.TrimRight(*baseURLPtr, "/"),
	}
	return cfg
}
