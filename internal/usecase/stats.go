package usecase

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
)

func (uc *useCase) GetStats(ctx context.Context) (models.Stats, error) {
	return uc.repository.GetStats(ctx)
}
