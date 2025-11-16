package repository

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/repository/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	var db *pgxpool.Pool
	if cfg.DataBaseDSN != "" {
		logger.Debugf("making new db connection pool with dsn: %v", cfg.DataBaseDSN)

		poolConfig, err := pgxpool.ParseConfig(cfg.DataBaseDSN)
		if err != nil {
			logger.WithError(err).Fatal("unable to parse DATABASE_URL")
		}

		db, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			logger.WithError(err).Fatal("unable to create connection pool")
		}

		m, err := migrate.New("file://migrations", cfg.DataBaseDSN)
		if err != nil {
			logger.WithError(err).Fatal("failed to create migrate instance")
		}

		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			logger.WithError(err).Fatal("failed to migrate")
		}
		logger.Info("db successfully migrated")
	}

	return &repository{
		storage: store,
		cfg:     cfg,
		logger:  logger,
		db:      db,
	}
}
