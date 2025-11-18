package storage

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
)

func (r *repository) GetShortIDsBatch(ctx context.Context, entities []*models.URLEntity) error {
	r.logger.Panic("GetShortIDsBatch batch is not implemented")
	return nil
}
