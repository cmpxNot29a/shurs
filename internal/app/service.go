package app

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/cmpxNot29a/shurs/internal/helper"
)

type ShortenerUseCase interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
	GetOriginalURL(ctx context.Context, id string) (string, error)
}

// ShortenerService инкапсулирует бизнес-логику сокращения URL.
type ShortenerService struct {
	storage  Storage
	idLength int
	attempts int
}

func NewShortenerService(storage Storage, idLength int, attempts int) *ShortenerService {

	return &ShortenerService{
		storage:  storage,
		idLength: idLength,
		attempts: attempts,
	}
}

// CreateShortURL генерирует уникальный ID, сохраняет URL и возвращает ID.
func (s *ShortenerService) CreateShortURL(ctx context.Context, originalURL string) (string, error) {

	shortID, err := s.genUnicID(ctx)
	if err != nil {
		return "", err
	}

	err = s.storage.Save(ctx, shortID, originalURL)
	if err == nil {
		return shortID, nil // Успешно сохранено
	}

	return "", fmt.Errorf("storage error during save: %w", err)
}

// GetOriginalURL получает оригинальный URL по ID.
func (s *ShortenerService) GetOriginalURL(ctx context.Context, id string) (string, error) {
	originalURL, err := s.storage.GetByID(ctx, id)
	if err != nil && !errors.Is(err, ErrNotFound) {
		log.Printf("ERROR (Service): Failed to get URL by ID %s: %v", id, err)
		return "", fmt.Errorf("storage error during get: %w", err)
	}
	return originalURL, err
}

func (s *ShortenerService) genUnicID(ctx context.Context) (string, error) {
	for range s.attempts {

		randomIDBytes, err := helper.GenerateRandomBase62(s.idLength)

		if err != nil {
			continue
		}

		randomID := string(randomIDBytes)
		exists, err := s.storage.Exists(ctx, randomID)

		if err != nil {
			log.Printf("ERROR (genUnicID): Failed to check existence for ID %s: %v", randomID, err)
			return "", fmt.Errorf("failed to check ID existence: %w", err)
		}
		if !exists {
			return randomID, nil
		}
		log.Printf("WARN: Handler: Collision detected for ID: %s. Retrying...", randomID)
	}
	return "", fmt.Errorf("failed to generate unique ID after %d attempts", s.attempts)
}
