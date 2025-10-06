package repository

import (
	"github.com/georgg2003/shortener/internal/models"
)

func (r repository) NewShortURL(url string, shortID string) {
	r.storage[shortID] = models.ShortURL{
		ShortID: shortID,
		LongURL: url,
	}
}
