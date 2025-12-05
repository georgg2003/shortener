package storage

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	repo "github.com/georgg2003/shortener/internal/repository"
)

func (r *repository) GetShortIDsBatch(ctx context.Context, entities []*models.URLEntity) (map[string]string, error) {
	r.logger.Fatal("GetShortIDsBatch batch is not implemented")
	return nil, repo.ErrNotImplemented
}
