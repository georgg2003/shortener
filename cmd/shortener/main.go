package main

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/georgg2003/shortener/api"
	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/delivery"
	"github.com/georgg2003/shortener/internal/delivery/grpcapi"
	"github.com/georgg2003/shortener/internal/pkg/interceptors"
	audit_repo "github.com/georgg2003/shortener/internal/repository/audit"
	audit_service "github.com/georgg2003/shortener/internal/repository/audit/service"
	audit_storage "github.com/georgg2003/shortener/internal/repository/audit/storage"
	audit_stub "github.com/georgg2003/shortener/internal/repository/audit/stub"
	"github.com/georgg2003/shortener/internal/repository/db"
	"github.com/georgg2003/shortener/internal/repository/db/postgres"
	"github.com/georgg2003/shortener/internal/repository/db/storage"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/georgg2003/shortener/internal/usecase/audit"
	"github.com/georgg2003/shortener/pkg/utils"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func newAuditRepos(cfg *config.Config, logger *logrus.Entry) []audit_repo.AuditRepository {
	repos := make([]audit_repo.AuditRepository, 0)
	if cfg.AuditFile != "" {
		repos = append(repos, audit_storage.New(logger, cfg.AuditFile))
	}
	if cfg.AuditURL != "" {
		repos = append(repos, audit_service.New(
			logger,
			cfg.AuditURL,
			audit_service.WithRetryCount(5), // TODO add to config
			audit_service.WithRetryWaitTime(time.Second),
		))
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

func listen(server *http.Server, logger *logrus.Logger) func() error {
	return func() error {
		logger.Infof("Listening on %v", server.Addr)
		var err error
		if server.TLSConfig != nil {
			err = server.ListenAndServeTLS("", "")
		} else {
			err = server.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return utils.ErrWrap(err, "http server failed")
		}
		return nil
	}
}

func listenGRPC(cfg *config.Config, server *grpc.Server, logger *logrus.Logger) func() error {
	return func() error {
		logger.Infof("Listening tcp for grpc on %v", cfg.GRPCListenAddr)
		listen, err := net.Listen("tcp", cfg.GRPCListenAddr)
		if err != nil {
			return utils.ErrWrap(err, "failed to init tcp listener")
		}
		reflection.Register(server)
		if err := server.Serve(listen); err != nil {
			return utils.ErrWrap(err, "failed to serve grpc server")
		}

		return nil
	}
}

func listenGRPCShutdown(ctx context.Context, server *grpc.Server, logger *logrus.Logger) func() error {
	return func() error {
		<-ctx.Done()
		logger.Infof("shutting down server")
		server.GracefulStop()
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

func printBuildInfo(logger logrus.FieldLogger) {
	logger.Infof("Build version: %s", buildVersion)
	logger.Infof("Build date: %s", buildDate)
	logger.Infof("Build commit: %s", buildCommit)
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	conf, err := config.Load()
	if err != nil {
		logger.WithError(err).Fatal("failed to load config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	repo := newRepo(ctx, conf, logger)
	usecase := usecase.New(repo, conf, logger)

	auditLogger := logger.WithField("subsystem", "audit")
	auditRepos := newAuditRepos(conf, auditLogger)
	observer := audit.NewAuditObserver(auditRepos, auditLogger)

	// Удаляет обсервер после завершения main
	defer usecase.Observe("audit", observer)()

	delivery := delivery.New(usecase, logger, conf)
	grpcServer := grpcapi.NewShortenerServer(usecase)

	r := delivery.GetNewRouter()

	var tlsConfig *tls.Config
	if conf.EnableHTTPS {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("cache-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("localhost"),
		}
		tlsConfig = manager.TLSConfig()
	}

	server := &http.Server{Addr: conf.ListenAddr, Handler: r, TLSConfig: tlsConfig}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(listen(server, logger))
	g.Go(listenShutdown(ctx, server, logger))

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RequestIDInterceptor,
			interceptors.NewLoggingInterceptor(logger),
			interceptors.NewAuthInterceptor(conf, logger, usecase),
		),
		grpc.Creds(insecure.NewCredentials()),
	)
	api.RegisterShortenerServiceServer(s, grpcServer)
	g.Go(listenGRPC(conf, s, logger))
	g.Go(listenGRPCShutdown(ctx, s, logger))

	if conf.DebugAddr != "" {
		debugServer := &http.Server{Addr: conf.DebugAddr, TLSConfig: tlsConfig}
		g.Go(listen(debugServer, logger))
		g.Go(listenShutdown(ctx, debugServer, logger))
	}

	printBuildInfo(logger)

	if err := g.Wait(); err != nil {
		logger.WithError(err).Fatal("application stopped with error")
	}
}
