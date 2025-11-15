package usecase

import "context"

func (uc *useCase) ProcessShortURL(ctx context.Context, id string) (string, error) {
	longURL, err := uc.repository.GetLongURL(ctx, id)
	return longURL, err
}
