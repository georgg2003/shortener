package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/delivery"
	audit_repo "github.com/georgg2003/shortener/internal/repository/audit"
	audit_service "github.com/georgg2003/shortener/internal/repository/audit/service"
	audit_storage "github.com/georgg2003/shortener/internal/repository/audit/storage"
	audit_stub "github.com/georgg2003/shortener/internal/repository/audit/stub"
	"github.com/georgg2003/shortener/internal/repository/db"
	"github.com/georgg2003/shortener/internal/repository/db/postgres"
	"github.com/georgg2003/shortener/internal/repository/db/storage"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/georgg2003/shortener/internal/usecase/audit"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func newAuditRepos(cfg *config.Config, logger *logrus.Entry) []audit_repo.AuditRepository {
	repos := make([]audit_repo.AuditRepository, 0)
	if cfg.AuditFile != "" {
		repos = append(repos, audit_storage.New(logger, cfg.AuditFile))
	}
	if cfg.AuditURL != "" {
		repos = append(repos, audit_service.New(logger, cfg.AuditURL))
	}
	if cfg.AuditStubEnabled {
		repos = append(repos, audit_stub.New(logger))
	}
	return repos
}

func newRepo(ctx context.Context, conf *config.Config, logger *logrus.Logger) db.Repository {
	if conf.DataBaseDSN != "" {
		return postgres.New(ctx, conf, logger)
	} else {
		return storage.New(ctx, conf, logger)
	}
}

func newConfig(logger *logrus.Logger) *config.Config {
	conf := config.New()
	if err := conf.ReadFromEnv(); err != nil {
		logger.WithError(err).Fatal("failed to read config from env")
	}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	if err := conf.ReadFromFlags(fs); err != nil {
		logger.WithError(err).Fatal("failed to read config from flags")
	}
	fs.Parse(os.Args[1:])

	return conf
}

func listen(server *http.Server, logger *logrus.Logger) func() error {
	return func() error {
		logger.Infof("Listening on %v", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Error("server failed")
			return err
		}
		return nil
	}
}

func listenShutdown(ctx context.Context, server *http.Server, logger *logrus.Logger) func() error {
	return func() error {
		<-ctx.Done()
		logger.Infof("shutting down server")
		return server.Shutdown(context.Background())
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	conf := newConfig(logger)
	repo := newRepo(ctx, conf, logger)
	usecase := usecase.New(repo, conf, logger)

	auditLogger := logger.WithField("subsystem", "audit")
	auditRepos := newAuditRepos(conf, auditLogger)
	observer := audit.NewAuditObserver(auditRepos)
	defer usecase.Observe("audit", observer)()

	delivery := delivery.New(usecase, logger, conf)

	r := delivery.GetNewRouter()

	server := &http.Server{Addr: conf.ListenAddr, Handler: r}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(listen(server, logger))
	g.Go(listenShutdown(ctx, server, logger))

	if conf.DebugAddr != "" {
		debugServer := &http.Server{Addr: conf.DebugAddr}
		g.Go(listen(debugServer, logger))
		g.Go(listenShutdown(ctx, debugServer, logger))
	}

	if err := g.Wait(); err != nil {
		logger.WithError(err).Fatal("application stopped with error")
	}
}
