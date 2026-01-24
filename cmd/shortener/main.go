package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/delivery"
	"github.com/georgg2003/shortener/internal/repository/audit_repo"
	audit_stup "github.com/georgg2003/shortener/internal/repository/audit_repo/stub"
	"github.com/georgg2003/shortener/internal/repository/db"
	"github.com/georgg2003/shortener/internal/repository/db/postgres"
	"github.com/georgg2003/shortener/internal/repository/db/storage"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/georgg2003/shortener/internal/usecase/audit"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func newAuditRepo(cfg *config.Config, logger *logrus.Entry) audit_repo.AuditRepository {
	return audit_stup.New(logger)
}

func newRepo(ctx context.Context, conf *config.Config, logger *logrus.Logger) db.Repository {
	if conf.DataBaseDSN != "" {
		return postgres.New(ctx, conf, logger)
	} else {
		return storage.New(ctx, conf, logger)
	}
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	conf := config.New()
	if err := conf.ReadFromEnv(); err != nil {
		logger.WithError(err).Fatal("failed to read config from env")
	}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	if err := conf.ReadFromFlags(fs); err != nil {
		logger.WithError(err).Fatal("failed to read config from flags")
	}
	fs.Parse(os.Args[1:])

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	repo := newRepo(ctx, conf, logger)
	usecase := usecase.New(repo, conf, logger)

	auditLogger := logger.WithField("subsystem", "audit")
	auditRepo := newAuditRepo(conf, auditLogger)
	observer := audit.NewAuditObserver(auditRepo)
	defer usecase.Observe("audit", observer)()

	delivery := delivery.New(usecase, logger, conf)

	r := delivery.GetNewRouter()

	server := &http.Server{Addr: conf.ListenAddr, Handler: r}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		logger.Infof("Listening on %v", conf.ListenAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Error("server failed")
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		logger.Infof("shutting down server")
		return server.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		logger.WithError(err).Error("application stopped with error")
	}
}
