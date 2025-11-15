package repository

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
)

func (r *repository) NewShortURL(ctx context.Context, url string, shortURL string) {
	r.storage.Store(shortURL, models.ShortURL{
		ShortURL: shortURL,
		LongURL:  url,
	})
}
