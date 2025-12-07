package storage

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	repo "github.com/georgg2003/shortener/internal/repository"
)

func (r *repository) DeleteUserURLs(ctx context.Context, tasks []models.DeleteUserURLsTask) error {
	r.logger.Fatal("GetShortIDsBatch batch is not implemented")
	return repo.ErrNotImplemented
}
