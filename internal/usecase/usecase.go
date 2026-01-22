package usecase

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/models"
	"github.com/sirupsen/logrus"
)

var errUserNotFound = errors.New("user is not found")

//go:generate mockgen -destination ./mock/mock.go -package mock . Repository
type Repository interface {
	NewShortURL(ctx context.Context, url string, shortID string) error
	GetLongURL(ctx context.Context, shortID string) (string, bool, error)
	GetShortID(ctx context.Context, originalURL string) (string, error)
	GetShortIDsBatch(ctx context.Context, entities []*models.URLEntity) (map[string]string, error)
	Ping(ctx context.Context) error
	NewShortURLBatch(ctx context.Context, entities []*models.URLEntity) error
	NewUser(ctx context.Context) (int64, error)
	GetUserURLs(ctx context.Context, userID int64) ([]models.URLEntity, error)
	DeleteUserURLs(ctx context.Context, tasks []models.DeleteUserURLsTask) error
}

type useCase struct {
	repository       Repository
	config           *config.Config
	deleteUserURLsCh chan models.DeleteUserURLsTask
	logger           *logrus.Logger
	observersByID    ObserversByID
}

func New(repo Repository, conf *config.Config, logger *logrus.Logger) *useCase {
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
