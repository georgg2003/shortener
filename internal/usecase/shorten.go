package usecase

import (
	"fmt"

	"github.com/georgg2003/shortener/pkg/utils"
)

const shortIDLength = 7

func (uc useCase) NewShortURL(url string) string {
	shortID := utils.RandomBase62(shortIDLength)
	uc.repository.NewShortURL(url, shortID)
	shortURL := fmt.Sprintf("%v/%v", uc.config.BaseURL, shortID)
	return shortURL
}
