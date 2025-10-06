package usecase

import "fmt"

func (s useCase) NewShortUrl(url string) string {
	// TODO write url to db
	shortUrlID := "abcasdds"
	s.repository.NewUrl(url, shortUrlID)
	shortUrl := fmt.Sprintf("http://localhost:8080/%v", shortUrlID)
	return shortUrl
}
