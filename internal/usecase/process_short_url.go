package usecase

import "errors"

func (s useCase) ProcessShortURL(id string) (string, error) {
	// TODO write url to db
	if id == "abcasdds" {
		return "https://practicum.yandex.ru/", nil
	}
	return "", errors.New("link not found")
}
