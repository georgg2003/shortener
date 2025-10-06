package usecase

import (
	"crypto/rand"
	"fmt"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const shortUrlIDLength = 7
const serviceURL = "http://localhost:8080"

func randomBase62(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphabet[int(b)%len(alphabet)]
	}
	return string(bytes)
}

func (s useCase) NewShortUrl(url string) string {
	shortUrlID := randomBase62(shortUrlIDLength)
	s.repository.NewShortUrl(url, shortUrlID)
	shortUrl := fmt.Sprintf("%v/%v/", serviceURL, shortUrlID)
	return shortUrl
}
