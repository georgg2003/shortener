package usecase

import "context"

func (uc *useCase) ProcessShortURL(ctx context.Context, id string) (string, bool, error) {
	longURL, isDeleted, err := uc.repository.GetLongURL(ctx, id)
	return longURL, isDeleted, err
}
