package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cmpxNot29a/shurs/internal/config"
	httpDelivery "github.com/cmpxNot29a/shurs/internal/delivery/http/handler"
	mw "github.com/cmpxNot29a/shurs/internal/delivery/http/middleware"
	"github.com/cmpxNot29a/shurs/internal/repository/memory"
	"github.com/cmpxNot29a/shurs/internal/service/idgenerator"
	shortenerUseCaseImpl "github.com/cmpxNot29a/shurs/internal/usecase/shortener"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func App(conf *config.Config) error {

	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatalf("FATAL: Can't initialize zap logger: %v", err)
	}

	defer func() {
		if errSync := logger.Sync(); errSync != nil {
			log.Printf("WARN: failed to sync zap logger: %v", errSync)
		}
	}()

	logger.Info("Logger initialized successfully")

	idLength := conf.IDLength
	attempts := conf.Attempts

	if idLength <= 0 {
		log.Printf("WARN (App): Invalid ID Length (%d) from config, using default %d", idLength, config.DefaultIDLength)
		idLength = config.DefaultIDLength // Используем константу из config
	}
	if attempts <= 0 {
		log.Printf("WARN (App): Invalid Attempts (%d) from config, using default %d", attempts, config.DefaultAttempts)
		attempts = config.DefaultAttempts // Используем константу из config
	}
	repo := memory.NewStorage()
	idGen := idgenerator.NewBase62Generator(idLength)
	shortenerUC := shortenerUseCaseImpl.NewUseCase(repo, idGen, attempts)

	handler := httpDelivery.NewHandler(shortenerUC, conf.BaseURL)

	r := chi.NewRouter()

	r.Use(mw.LoggingMiddleware(logger))

	r.Post("/", mw.ValidateURL(http.HandlerFunc(handler.Create)).ServeHTTP)
	r.Get("/{id}", mw.ValidateID(idLength)(http.HandlerFunc(handler.Redirect)).ServeHTTP)

	logger.Info("Starting server", // Используем zap
		zap.String("ServerAddress", conf.ServerAddress),
		zap.String("BaseURL", conf.BaseURL),
	)

	err = http.ListenAndServe(conf.ServerAddress, r)
	if err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}
	return nil
}
