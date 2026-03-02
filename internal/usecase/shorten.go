package usecase

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/contextlib"
	"github.com/georgg2003/shortener/pkg/postgres"
	"github.com/georgg2003/shortener/pkg/utils"
)

const shortIDLength = 8

var errCreatingNewShortURLFailed = errors.New("creation of a new short url was failed")
var ErrURLEntityAlreadyExists = errors.New("conflict, url is already shortened")

func (uc *useCase) shortURLFromID(shortID string) string {
	return fmt.Sprintf("%v/%v", uc.config.BaseURL, shortID)
}

var errFailedToGetShortID = errors.New("failed to get short id")

func (uc *useCase) NewShortURL(ctx context.Context, url string) (string, error) {
	shortID, err := uc.base62Gen.Generate(shortIDLength)
	if err != nil {
		return "", utils.ErrWrap(err, "failed to generate short ID")
	}
	userID, _ := contextlib.GetUserID(ctx)
	err = uc.repository.NewShortURL(ctx, url, shortID)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			shortID, err = uc.repository.GetShortID(ctx, url)
			if err != nil {
				return "", utils.ErrWrap(err, errFailedToGetShortID.Error())
			}
			err = ErrURLEntityAlreadyExists
		} else {
			return "", utils.ErrWrap(err, errCreatingNewShortURLFailed.Error())
		}
	}
	shortURL := uc.shortURLFromID(shortID)
	if uc.observersByID != nil && err == nil {
		for observer := range maps.Values(uc.observersByID) {
			observer.OnNewURL(ObserverEvent{
				Time:        time.Now(),
				UserID:      userID,
				OriginalURL: url,
			})
		}
	}
	return shortURL, err
}

var errFailedToGetShortIDsBatch = errors.New("failed to get short ids batch")

func (uc *useCase) NewShortURLBatch(ctx context.Context, entities []*models.URLEntity) error {
	for _, v := range entities {
		shortID, err := uc.base62Gen.Generate(shortIDLength)
		if err != nil {
			return utils.ErrWrap(err, "failed to generate short ID")
		}
		v.ShortID = shortID
		v.ShortURL = uc.shortURLFromID(shortID)
	}

	err := uc.repository.NewShortURLBatch(ctx, entities)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			existingIDs, getErr := uc.repository.GetShortIDsBatch(ctx, entities)
			if getErr != nil {
				return utils.ErrWrap(getErr, errFailedToGetShortIDsBatch.Error())
			}
			for _, v := range entities {
				if existingID, ok := existingIDs[v.OriginalURL]; ok {
					v.ShortID = existingID
					v.ShortURL = uc.shortURLFromID(existingID)
				}
			}
			return ErrURLEntityAlreadyExists
		}
		return utils.ErrWrap(err, errCreatingNewShortURLFailed.Error())
	}

	return nil
}
