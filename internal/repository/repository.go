package repository

import (
	"context"
	"sync"
)

type Repository interface {
	NewShortURL(ctx context.Context, url string, shortID string)
	GetLongURL(ctx context.Context, shortID string) (string, error)
}

type repository struct {
	storage *sync.Map
}

func New() Repository {
	storage := sync.Map{}

	return repository{
		storage: &storage,
	}
}
