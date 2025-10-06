package repository

import "github.com/georgg2003/shortener/internal/models"

type Repository interface {
	NewShortURL(url string, shortID string)
	GetLongURL(shortID string) (string, error)
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
