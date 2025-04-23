package shortener

import (
	"context"
	"errors"
	"fmt"
	"log"

	repo "github.com/cmpxNot29a/shurs/internal/domain/repository"
	"github.com/cmpxNot29a/shurs/internal/domain/service"
	domainUC "github.com/cmpxNot29a/shurs/internal/domain/usecase"
)

type useCase struct {
	repo     repo.ShortLinkRepository
	idGen    service.IDGenerator
	attempts int
}

// NewUseCase создает новый экземпляр useCase, возвращая интерфейс из domain.
func NewUseCase(repo repo.ShortLinkRepository, idGen service.IDGenerator, attempts int) domainUC.ShortenerUseCase {
	return &useCase{
		repo:     repo,
		idGen:    idGen,
		attempts: attempts,
	}
}

// genUniqueID - неэкспортируемый метод для генерации уникального ID.
func (uc *useCase) genUniqueID(ctx context.Context) (string, error) {
	for range uc.attempts {
		shortID, err := uc.idGen.Generate() // Используем интерфейс генератора
		if err != nil {
			continue
		}

		exists, err := uc.repo.Exists(ctx, shortID)
		if err != nil {
			log.Printf("ERROR (genUniqueID): Failed to check existence for ID %s: %v", shortID, err)
			return "", fmt.Errorf("failed to check ID existence: %w", err)
		}

		if !exists {
			return shortID, nil
		}

		log.Printf("WARN (genUniqueID): Collision detected for ID %s. Retrying...", shortID)
	}
	return "", fmt.Errorf("failed to generate unique ID after %d attempts", uc.attempts)
}

// CreateShortURL реализует метод интерфейса domainUC.ShortenerUseCase.
func (uc *useCase) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	shortID, err := uc.genUniqueID(ctx)
	if err != nil {
		log.Printf("ERROR (UseCase): Failed to generate unique ID: %v", err)
		return "", fmt.Errorf("could not generate unique short ID: %w", err)
	}

	err = uc.repo.Save(ctx, shortID, originalURL)
	if err != nil {
		if errors.Is(err, repo.ErrConflict) {
			log.Printf("ERROR (UseCase): Conflict during save for supposedly unique ID %s", shortID)
			return "", fmt.Errorf("unexpected conflict during save for ID %s", shortID)
		}
		log.Printf("ERROR (UseCase): Failed to save URL for ID %s: %v", shortID, err)
		return "", fmt.Errorf("storage error during save: %w", err)
	}
	return shortID, nil
}

// GetOriginalURL реализует метод интерфейса domainUC.ShortenerUseCase.
func (uc *useCase) GetOriginalURL(ctx context.Context, shortID string) (string, error) {
	originalURL, err := uc.repo.GetByID(ctx, shortID)
	if err != nil && !errors.Is(err, repo.ErrNotFound) {
		log.Printf("ERROR (UseCase): Failed to get URL by ID %s: %v", shortID, err)
		return "", fmt.Errorf("storage error during get: %w", err)
	}
	return originalURL, err
}
