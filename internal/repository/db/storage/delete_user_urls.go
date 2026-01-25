package storage

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/internal/repository/db"
)

func (r *repository) DeleteUserURLs(ctx context.Context, tasks []models.DeleteUserURLsTask) error {
	r.logger.Fatal("GetShortIDsBatch batch is not implemented")
	return db.ErrNotImplemented
}
