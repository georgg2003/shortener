package usecase

import (
	"context"
	"maps"
	"time"

	"github.com/georgg2003/shortener/pkg/contextlib"
)

func (uc *useCase) ProcessShortURL(ctx context.Context, id string) (string, bool, error) {
	longURL, isDeleted, err := uc.repository.GetLongURL(ctx, id)
	userID, _ := contextlib.GetUserID(ctx)

	if uc.observersByID != nil {
		for observer := range maps.Values(uc.observersByID) {
			observer.OnNewURL(ObserverEvent{
				Time:        time.Now(),
				UserID:      userID,
				OriginalURL: longURL,
			})
		}
	}

	return longURL, isDeleted, err
}
