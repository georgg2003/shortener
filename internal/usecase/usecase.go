package usecase

import (
	"context"

	"github.com/georgg2003/shortener/internal/config"
)

//go:generate mockgen -destination ./mock/mock.go -package mock . Repository
type Repository interface {
	NewShortURL(ctx context.Context, url string, shortID string)
	GetLongURL(ctx context.Context, shortID string) (string, error)
	Ping(ctx context.Context) error
}

type useCase struct {
	repository Repository
	config     *config.Config
}

func New(repo Repository, conf *config.Config) *useCase {
	return &useCase{
		repository: repo,
		config:     conf,
	}
}
