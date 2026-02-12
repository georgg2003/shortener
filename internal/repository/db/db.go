package db

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/models"
)

var ErrNotFound = errors.New("value not found")
var ErrNotImplemented = errors.New("method not implemented")

//go:generate go tool mockgen -destination ./mock/mock.go -package mock . Repository
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
