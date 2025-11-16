package repository

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/repository/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var errFailedToAcquireConnection = errors.New("failed to acquire db conn")
var errFailedToScan = errors.New("failed to scan a row")

type repository struct {
	storage *storage.SyncMapStorage
	db      *pgxpool.Pool
	cfg     *config.Config
	logger  *logrus.Logger
}

func New(
	ctx context.Context,
	cfg *config.Config,
	logger *logrus.Logger,
) *repository {
	store := storage.New(
		&storage.SyncStorageConfig{
			FileStoragePath: cfg.FileStoragePath,
		},
		logger,
	)

	logger.Debugf("making new db connection pool with dsn: %v", cfg.DataBaseDSN)

	poolConfig, err := pgxpool.ParseConfig(cfg.DataBaseDSN)
	if err != nil {
		logger.Fatal("Unable to parse DATABASE_URL:", err)
	}

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Fatal("Unable to create connection pool:", err)
	}

	return &repository{
		storage: store,
		cfg:     cfg,
		logger:  logger,
		db:      db,
	}
}
