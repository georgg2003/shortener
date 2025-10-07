package usecase

import (
	repo "github.com/georgg2003/shortener/internal/repository"
)

type UseCase interface {
	NewShortURL(url string, baseURL string) string
	ProcessShortURL(id string) (string, error)
}

type useCase struct {
	repository repo.Repository
}

func New(repo repo.Repository) UseCase {
	return useCase{
		repository: repo,
	}
}
