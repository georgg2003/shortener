package usecase

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/internal/repository/db"
	"github.com/sirupsen/logrus"
)

var errUserNotFound = errors.New("user is not found")

type UseCase interface {
	NewShortURL(ctx context.Context, url string) (string, error)
	NewShortURLBatch(ctx context.Context, entities []*models.URLEntity) error
	ProcessShortURL(ctx context.Context, id string) (string, bool, error)
	Ping(ctx context.Context) error
	NewUser(ctx context.Context) (int64, error)
	GetUserURLs(ctx context.Context) ([]models.URLEntity, error)
	DeleteUserURLs(ctx context.Context, urls []string) error
}

type useCase struct {
	repository       db.Repository
	config           *config.Config
	deleteUserURLsCh chan models.DeleteUserURLsTask
	logger           *logrus.Logger
	observersByID    observersByID
}

func New(repo db.Repository, conf *config.Config, logger *logrus.Logger) *useCase {
	deleteUserURLsCh := make(chan models.DeleteUserURLsTask, 1024)

	uc := useCase{
		repository:       repo,
		config:           conf,
		deleteUserURLsCh: deleteUserURLsCh,
		logger:           logger,
	}

	go uc.deleteUserURLsWorker()

	return &uc
}
