// Модуль, описывающий REST API ручки сервиса сокращения URL.
package delivery

import (
	"io"
	"net/http"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/georgg2003/shortener/pkg/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

const internalErrorText = "Internal Error"

type Delivery interface {
	GetNewRouter() chi.Router
	ShortenURL(w http.ResponseWriter, r *http.Request)
	ProcessShortURL(w http.ResponseWriter, r *http.Request)
	APIShortenURL(w http.ResponseWriter, r *http.Request)
	APIShortenURLBatch(w http.ResponseWriter, r *http.Request)
	GetUserURLs(w http.ResponseWriter, r *http.Request)
	DeleteUserURLs(w http.ResponseWriter, r *http.Request)
	Ping(w http.ResponseWriter, r *http.Request)
}

type delivery struct {
	usecase usecase.UseCase
	logger  logrus.FieldLogger
	cfg     *config.Config
}

func New(
	usecase usecase.UseCase,
	logger logrus.FieldLogger,
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
		middlewares.NewGzipCompressionMiddleware(d.logger),
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

func (d *delivery) safeClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		d.logger.WithError(err).Error("failed to close")
	}
}
