package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/utils"
)

const shortIDLength = 8

var errCreatingNewShortURLFailed = errors.New("creation of a new short url was failed")

func (uc *useCase) NewShortURL(ctx context.Context, url string) (string, error) {
	shortID := utils.RandomBase62(shortIDLength)
	err := uc.repository.NewShortURL(ctx, url, shortID)
	if err != nil {
		return "", errors.Join(err, errCreatingNewShortURLFailed)
	}
	shortURL := fmt.Sprintf("%v/%v", uc.config.BaseURL, shortID)
	return shortURL, nil
}

func (uc *useCase) NewShortURLBatch(ctx context.Context, entities []*models.URLEntity) error {
	for _, v := range entities {
		shortID := utils.RandomBase62(shortIDLength)
		v.ShortID = shortID
		v.ShortURL = fmt.Sprintf("%v/%v", uc.config.BaseURL, shortID)
	}

	log.Println(entities)

	err := uc.repository.NewShortURLBatch(ctx, entities)
	if err != nil {
		return errors.Join(err, errCreatingNewShortURLFailed)
	}

	return nil
}
