package delivery

import (
	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/georgg2003/shortener/pkg/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

const internalErrorText = "Internal Error"

type Delivery interface {
	GetNewRouter() chi.Router
}

type delivery struct {
	usecase usecase.UseCase
	logger  *logrus.Logger
	cfg     *config.Config
}

func New(
	usecase usecase.UseCase,
	logger *logrus.Logger,
	cfg *config.Config,
) Delivery {
	return &delivery{
		usecase: usecase,
		logger:  logger,
		cfg:     cfg,
	}
}

func (d *delivery) GetNewRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(
		middlewares.NewAccessLogMiddleware(d.logger),
		middlewares.NewGzipCompressionMiddleware(),
	)

	if d.cfg.DataBaseDSN != "" {
		r.Use(middlewares.NewSimpleAuthMiddleware(d.cfg, d.logger, d.usecase))
	}

	r.Post("/api/shorten", d.APIShortenURL)
	r.Post("/api/shorten/batch", d.APIShortenURLBatch)
	r.Get("/api/user/urls", d.GetUserURLs)
	r.Delete("/api/user/urls", d.DeleteUserURLs)

	r.Post("/", d.ShortenURL)
	r.Get("/{id}", d.ProcessShortURL)
	r.Get("/ping", d.Ping)

	return r
}
