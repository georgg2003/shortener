package repository

import (
	"errors"
)

func (r repository) GetLongURL(shortID string) (string, error) {
	shortURLModel, ok := r.storage[shortID]
	if !ok {
		return "", errors.New("value not found")
	}
	return shortURLModel.LongURL, nil
}
