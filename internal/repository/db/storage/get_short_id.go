package storage

import (
	"context"

	"github.com/georgg2003/shortener/internal/repository/db"
)

func (r *repository) GetShortID(ctx context.Context, originalURL string) (string, error) {
	r.logger.Fatal("GetShortID is not implemented")
	return "", db.ErrNotImplemented
}
