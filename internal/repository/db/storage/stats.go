package storage

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/internal/repository/db"
)

func (r *repository) GetStats(ctx context.Context) (models.Stats, error) {
	r.logger.Fatal("GetStats method is not implemented")
	return models.Stats{}, db.ErrNotImplemented
}
