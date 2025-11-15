package repository

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/models"
)

var errNotFound = errors.New("value not found")
var errWrongStorageType = errors.New("wrong type in storage")

func (r *repository) GetLongURL(ctx context.Context, shortID string) (string, error) {
	storageValue, ok := r.storage.Load(shortID)
	if !ok {
		return "", errNotFound
	}

	shortURLModel, ok := storageValue.(models.ShortURL)
	if !ok {
		return "", errWrongStorageType
	}

	return shortURLModel.LongURL, nil
}
