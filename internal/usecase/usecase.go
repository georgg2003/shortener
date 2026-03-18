package usecase

import (
	"context"
	"errors"
	"slices"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/internal/repository/db"
	"github.com/georgg2003/shortener/pkg/utils"
	"github.com/sirupsen/logrus"
)

var errUserNotFound = errors.New("user is not found")

//go:generate go tool mockgen -destination ./mock/mock.go -package mock . UseCase
type UseCase interface {
	NewShortURL(ctx context.Context, url string) (shortURL string, err error)
	NewShortURLBatch(ctx context.Context, entities []*models.URLEntity) error
	ProcessShortURL(ctx context.Context, id string) (originalURL string, isDeleted bool, err error)
	Ping(ctx context.Context) error
	NewUser(ctx context.Context) (int64, error)
	GetUserURLs(ctx context.Context) ([]models.URLEntity, error)
	DeleteUserURLs(ctx context.Context, urls []string) error
	GetStats(ctx context.Context) (models.Stats, error)
}

type useCase struct {
	repository       db.Repository
	config           *config.Config
	deleteUserURLsCh chan models.DeleteUserURLsTask
	logger           *logrus.Logger
	observersByID    observersByID
	base62Gen        utils.Base62Generator
}

type UseCaseOption func(uc *useCase)

func WithBase62Generator(generator utils.Base62Generator) UseCaseOption {
	return func(uc *useCase) {
		uc.base62Gen = generator
	}
}

func New(repo db.Repository, conf *config.Config, logger *logrus.Logger, opts ...UseCaseOption) *useCase {
	deleteUserURLsCh := make(chan models.DeleteUserURLsTask, 1024)

	uc := useCase{
		repository:       repo,
		config:           conf,
		deleteUserURLsCh: deleteUserURLsCh,
		logger:           logger,
		base62Gen:        utils.CryptoBase62Generator{},
	}

	for opt := range slices.Values(opts) {
		opt(&uc)
	}

	go uc.deleteUserURLsWorker()

	return &uc
}
