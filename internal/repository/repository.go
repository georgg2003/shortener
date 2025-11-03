package repository

import (
	"context"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	NewShortURL(ctx context.Context, url string, shortID string)
	GetLongURL(ctx context.Context, shortID string) (string, error)
}

type repository struct {
	storage *storage
	cfg     *config.Config
	logger  *logrus.Logger
}

func New(
	cfg *config.Config,
	logger *logrus.Logger,
) Repository {
	store := NewStorage(
		cfg,
		logger,
	)

	repo := repository{
		storage: store,
		cfg:     cfg,
		logger:  logger,
	}

	return repo
}
