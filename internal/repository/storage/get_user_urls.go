package storage

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	repo "github.com/georgg2003/shortener/internal/repository"
)

func (r *repository) GetUserURLs(ctx context.Context, userID int64) ([]models.URLEntity, error) {
	r.logger.Fatal("GetShortIDsBatch batch is not implemented")
	return nil, repo.ErrNotImplemented
}
