package usecase

import (
	"context"

	"github.com/georgg2003/shortener/internal/config"
	repo "github.com/georgg2003/shortener/internal/repository"
)

type UseCase interface {
	NewShortURL(ctx context.Context, url string) string
	ProcessShortURL(ctx context.Context, id string) (string, error)
}

type useCase struct {
	repository repo.Repository
	config     *config.Config
}

func New(repo repo.Repository, conf *config.Config) UseCase {
	return useCase{
		repository: repo,
		config:     conf,
	}
}
