package storage

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/internal/repository/db"
)

var errWrongStorageType = errors.New("wrong type in storage")

func (r *repository) GetLongURL(ctx context.Context, shortID string) (string, bool, error) {
	storageValue, ok := r.storage.Load(shortID)
	if !ok {
		return "", false, db.ErrNotFound
	}

	shortURLModel, ok := storageValue.(models.ShortURL)
	if !ok {
		return "", false, errWrongStorageType
	}

	return shortURLModel.LongURL, false, nil
}
