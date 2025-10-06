package repository

import (
	"github.com/georgg2003/shortener/internal/models"
)

func (r repository) NewShortURL(url string, shortUrlID string) {
	r.storage[shortUrlID] = models.ShortURL{
		ShortUrlID: shortUrlID,
		LongURL:    url,
	}
}
