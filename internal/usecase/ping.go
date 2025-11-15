package usecase

import "context"

func (uc *useCase) Ping(ctx context.Context) error {
	return uc.repository.Ping(ctx)
}
