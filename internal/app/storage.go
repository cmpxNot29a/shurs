package app

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("short link not found")
	ErrConflict = errors.New("short link ID conflict or already exists")
)

type Storage interface {
	Save(ctx context.Context, id, originalURL string) error
	GetByID(ctx context.Context, id string) (originalURL string, err error)
	Exists(ctx context.Context, id string) (bool, error)
	Close() error
}
