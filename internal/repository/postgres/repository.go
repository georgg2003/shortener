package postgres

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var errFailedToAcquireConnection = errors.New("failed to acquire db conn")
var errFailedToBeginTransaction = errors.New("failed to begin transaction")
var errFailedToScan = errors.New("failed to scan a row")

type repository struct {
	db     *pgxpool.Pool
	cfg    *config.Config
	logger *logrus.Logger
}

func New(
	ctx context.Context,
	cfg *config.Config,
	logger *logrus.Logger,
) *repository {
	var db *pgxpool.Pool

	if cfg.DataBaseDSN == "" {
		logger.Fatalf("database dsn is empty")
	}
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

	return &repository{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}
