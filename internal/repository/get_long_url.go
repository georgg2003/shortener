package repository

import (
	"errors"
)

func (r repository) GetLongURL(shortUrlID string) (string, error) {
	shortURLModel, ok := r.storage[shortUrlID]
	if !ok {
		return "", errors.New("value not found")
	}
	return shortURLModel.LongURL, nil
}
