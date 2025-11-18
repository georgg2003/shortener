package storage

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
)

func (r *repository) GetShortIDsBatch(ctx context.Context, entities []*models.URLEntity) (map[string]string, error) {
	r.logger.Panic("GetShortIDsBatch batch is not implemented")
	return nil, nil
}
