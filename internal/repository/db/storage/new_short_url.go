package storage

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/models"
)

var errFailedToInsertNewShortURL = errors.New("failed to insert new short url")

func (r *repository) NewShortURL(ctx context.Context, url string, shortID string) error {
	r.storage.Store(shortID, models.ShortURL{
		ShortURL: shortID,
		LongURL:  url,
	})
	return nil
}
