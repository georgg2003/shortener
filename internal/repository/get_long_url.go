package repository

import (
	"errors"
)

func (r repository) GetLongUrl(shortUrlID string) (string, error) {
	shortURLModel, ok := r.storage[shortUrlID]
	if ok != true {
		return "", errors.New("value not found")
	}
	return shortURLModel.LongURL, nil
}
