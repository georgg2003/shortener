package storage

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/internal/repository/db"
)

func (r *repository) GetShortIDsBatch(ctx context.Context, entities []*models.URLEntity) (map[string]string, error) {
	r.logger.Fatal("GetShortIDsBatch batch is not implemented")
	return nil, db.ErrNotImplemented
}
