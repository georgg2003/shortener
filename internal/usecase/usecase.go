package usecase

import (
	"github.com/georgg2003/shortener/internal/config"
	repo "github.com/georgg2003/shortener/internal/repository"
)

type UseCase interface {
	NewShortURL(url string) string
	ProcessShortURL(id string) (string, error)
}

type useCase struct {
	repository repo.Repository
	config     config.Config
}

func New(repo repo.Repository, conf config.Config) UseCase {
	return useCase{
		repository: repo,
		config:     conf,
	}
}
