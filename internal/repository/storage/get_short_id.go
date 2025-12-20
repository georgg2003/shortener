package storage

import (
	"context"

	repo "github.com/georgg2003/shortener/internal/repository"
)

func (r *repository) GetShortID(ctx context.Context, originalURL string) (string, error) {
	r.logger.Fatal("GetShortID is not implemented")
	return "", repo.ErrNotImplemented
}
