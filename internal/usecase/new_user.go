package usecase

import "context"

func (uc *useCase) NewUser(ctx context.Context) (int64, error) {
	return uc.repository.NewUser(ctx)
}
