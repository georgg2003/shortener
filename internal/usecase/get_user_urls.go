package usecase

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/contextlib"
	"github.com/georgg2003/shortener/pkg/utils"
)

var errUserNotFound = errors.New("user is not found")

func (uc *useCase) GetUserURLs(ctx context.Context) ([]models.URLEntity, error) {
	userID, ok := contextlib.GetUserID(ctx)
	if !ok {
		return nil, errUserNotFound
	}

	urls, err := uc.repository.GetUserURLs(ctx, userID)
	if err != nil {
		return nil, utils.ErrWrap(err, "failed to get user urls")
	}

	for _, url := range urls {
		url.ShortURL = uc.shortURLFromID(url.ShortID)
	}

	return urls, err
}
