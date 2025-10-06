package usecase

import (
	"fmt"

	"github.com/georgg2003/shortener/pkg/utils"
)

const shortIDLength = 7
const serviceURL = "http://localhost:8080"

func (s useCase) NewShortURL(url string) string {
	shortID := utils.RandomBase62(shortIDLength)
	s.repository.NewShortURL(url, shortID)
	shortUrl := fmt.Sprintf("%v/%v", serviceURL, shortID)
	return shortUrl
}
