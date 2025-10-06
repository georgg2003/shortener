package repository

import (
	"github.com/georgg2003/shortener/internal/models"
)

func (r repository) NewShortUrl(url string, shortUrlID string) {
	r.storage[shortUrlID] = models.ShortURL{
		ShortUrlID: shortUrlID,
		LongURL:    url,
	}
}
