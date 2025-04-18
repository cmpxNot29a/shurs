package app

import (
	"context"
)

type InMemoryStorage struct {
	data map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

// Save реализует метод интерфейса Storage.
func (s *InMemoryStorage) Save(ctx context.Context, id, originalURL string) error {

	if _, exists := s.data[id]; exists {
		return ErrConflict
	}
	s.data[id] = originalURL
	return nil
}

// GetByID реализует метод интерфейса Storage.
func (s *InMemoryStorage) GetByID(ctx context.Context, id string) (string, error) {

	url, exists := s.data[id]
	if !exists {
		return "", ErrNotFound
	}
	return url, nil
}

// Exists реализует метод интерфейса Storage.
func (s *InMemoryStorage) Exists(ctx context.Context, id string) (bool, error) {

	_, exists := s.data[id]
	return exists, nil

}

// Close реализует метод интерфейса Storage.
func (s *InMemoryStorage) Close() error {
	return nil
}
