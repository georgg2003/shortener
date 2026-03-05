package storage

import (
	"context"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/pkg/storage"
	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type repository struct {
	storage *storage.SyncMapStorage
	cfg     *config.Config
	logger  *logrus.Logger
}

func New(
	ctx context.Context,
	cfg *config.Config,
	logger *logrus.Logger,
) *repository {
	store := storage.New(
		ctx,
		storage.SyncStorageConfig{
			FileStoragePath: cfg.FileStoragePath,
		},
		logger,
	)

	return &repository{
		storage: store,
		cfg:     cfg,
		logger:  logger,
	}
}
