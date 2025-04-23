package memory

import (
	"context"
	"github.com/cmpxNot29a/shurs/internal/domain/repository"
)

// Storage реализует repository.ShortLinkRepository для хранения в памяти.
type Storage struct {
	data map[string]string
}

// NewStorage создает новый экземпляр Storage.
func NewStorage() *Storage {
	return &Storage{data: make(map[string]string)}
}

// Save реализует метод интерфейса repository.ShortLinkRepository.
func (s *Storage) Save(ctx context.Context, id, originalURL string) error {
	if _, exists := s.data[id]; exists {
		return repository.ErrConflict
	}
	s.data[id] = originalURL
	return nil
}

// GetByID реализует метод интерфейса repository.ShortLinkRepository.
func (s *Storage) GetByID(ctx context.Context, id string) (string, error) {
	url, exists := s.data[id]
	if !exists {
		return "", repository.ErrNotFound
	}
	return url, nil
}

// Exists реализует метод интерфейса repository.ShortLinkRepository.
func (s *Storage) Exists(ctx context.Context, id string) (bool, error) {
	_, exists := s.data[id]
	return exists, nil
}

// Close реализует метод интерфейса repository.ShortLinkRepository.
func (s *Storage) Close() error {
	return nil
}
