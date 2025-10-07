package usecase

import (
	"fmt"

	"github.com/georgg2003/shortener/pkg/utils"
)

const shortIDLength = 7

func (s useCase) NewShortURL(url string, baseURL string) string {
	shortID := utils.RandomBase62(shortIDLength)
	s.repository.NewShortURL(url, shortID)
	shortURL := fmt.Sprintf("%v/%v", baseURL, shortID)
	return shortURL
}
