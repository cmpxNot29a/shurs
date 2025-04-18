package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cmpxNot29a/shurs/internal/config"
	"github.com/go-chi/chi/v5"
)

func App(conf *config.Config) error {

	idLength := conf.IDLength
	attempts := conf.Attempts

	storage := NewInMemoryStorage()

	if idLength <= 0 {
		log.Printf("WARN (App): Invalid ID Length (%d) from config, using default %d", idLength, config.DefaultIDLength)
		idLength = config.DefaultIDLength // Используем константу из config
	}
	if attempts <= 0 {
		log.Printf("WARN (App): Invalid Attempts (%d) from config, using default %d", attempts, config.DefaultAttempts)
		attempts = config.DefaultAttempts // Используем константу из config
	}

	var service ShortenerUseCase = NewShortenerService(storage, idLength, attempts)
	handler := NewHandler(service, conf.BaseURL)

	r := chi.NewRouter()

	idValidatorMiddleware := ValidateIDMiddleware(idLength)
	r.Post("/", ValidateURLMiddleware(http.HandlerFunc(handler.CreateShortURL)).ServeHTTP)
	r.Get("/{id}", idValidatorMiddleware(http.HandlerFunc(handler.Redirect)).ServeHTTP)

	log.Printf("INFO: Starting server on address %s", conf.ServerAddress)
	err := http.ListenAndServe(conf.ServerAddress, r)
	if err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}
	return nil
}
