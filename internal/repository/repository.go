package repository

import "github.com/georgg2003/shortener/internal/models"

type Repository interface {
	NewShortUrl(url string, shortUrlID string)
	GetLongUrl(shortUrlID string) (string, error)
}

type repository struct {
	storage map[string]models.ShortURL
}

func New() Repository {
	storage := make(map[string]models.ShortURL)

	return repository{
		storage: storage,
	}
}
