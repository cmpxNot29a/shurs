package usecase

import (
	"context"
)

type ShortenerUseCase interface {
	CreateShortURL(ctx context.Context, originalURL string) (shortID string, err error)
	GetOriginalURL(ctx context.Context, shortID string) (originalURL string, err error)
}
