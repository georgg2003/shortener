package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/postgres"
	"github.com/georgg2003/shortener/pkg/utils"
)

const shortIDLength = 8

var errCreatingNewShortURLFailed = errors.New("creation of a new short url was failed")
var ErrUrlEntityAlreadyExists = errors.New("conflict, url is already shortened")

func (uc *useCase) shortURLFromID(shortID string) string {
	return fmt.Sprintf("%v/%v", uc.config.BaseURL, shortID)
}

func (uc *useCase) NewShortURL(ctx context.Context, url string) (string, error) {
	shortID := utils.RandomBase62(shortIDLength)
	err := uc.repository.NewShortURL(ctx, url, shortID)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			shortID, err = uc.repository.GetShortID(ctx, url)
			if err == nil {
				err = ErrUrlEntityAlreadyExists
			}
		} else {
			return "", utils.ErrWrap(err, errCreatingNewShortURLFailed)
		}
	}
	shortURL := uc.shortURLFromID(shortID)
	return shortURL, err
}

func (uc *useCase) NewShortURLBatch(ctx context.Context, entities []*models.URLEntity) error {
	for _, v := range entities {
		shortID := utils.RandomBase62(shortIDLength)
		v.ShortID = shortID
		v.ShortURL = uc.shortURLFromID(shortID)
	}

	err := uc.repository.NewShortURLBatch(ctx, entities)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			err := uc.repository.GetShortIDsBatch(ctx, entities)
			if err == nil {
				err = ErrUrlEntityAlreadyExists
			}
			return err
		}
		return utils.ErrWrap(err, errCreatingNewShortURLFailed)
	}

	return nil
}
