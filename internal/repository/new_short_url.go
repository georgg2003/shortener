package repository

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
)

func (r repository) NewShortURL(ctx context.Context, url string, shortID string) {
	r.storage.Store(shortID, models.ShortURL{
		ShortID: shortID,
		LongURL: url,
	})
}
