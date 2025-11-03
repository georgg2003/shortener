package usecase

import (
	"context"
	"fmt"

	"github.com/georgg2003/shortener/pkg/utils"
)

const shortIDLength = 8

func (uc useCase) NewShortURL(ctx context.Context, url string) string {
	shortID := utils.RandomBase62(shortIDLength)
	uc.repository.NewShortURL(ctx, url, shortID)
	shortURL := fmt.Sprintf("%v/%v", uc.config.BaseURL, shortID)
	return shortURL
}
