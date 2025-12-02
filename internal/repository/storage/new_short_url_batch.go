package storage

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
)

func (r *repository) NewShortURLBatch(ctx context.Context, entities []*models.URLEntity) error {
	for _, v := range entities {
		r.storage.Store(v.ShortID, models.ShortURL{
			ShortURL: v.ShortID,
			LongURL:  v.OriginalURL,
		})
	}
	return nil
}
