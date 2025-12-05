package delivery

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

const internalErrorText = "Internal Error"

type Delivery interface {
	GetNewRouter() chi.Router
}

type UseCase interface {
	NewShortURL(ctx context.Context, url string) (string, error)
	NewShortURLBatch(ctx context.Context, entities []*models.URLEntity) error
	ProcessShortURL(ctx context.Context, id string) (string, error)
	Ping(ctx context.Context) error
	NewUser(ctx context.Context) (int64, error)
	GetUserURLs(ctx context.Context) ([]models.URLEntity, error)
}

type delivery struct {
	usecase UseCase
	logger  *logrus.Logger
}

func New(
	usecase UseCase,
	logger *logrus.Logger,
) Delivery {
	return &delivery{
		usecase: usecase,
		logger:  logger,
	}
}

func (d *delivery) GetNewRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(
		middlewares.NewAccessLogMiddleware(d.logger),
		middlewares.NewSimpleAuthMiddleware(d.logger, d.usecase),
		middlewares.NewGzipCompressionMiddleware(),
	)

	r.Post("/api/shorten", d.APIShortenURL)
	r.Post("/api/shorten/batch", d.APIShortenURLBatch)
	r.Get("/api/user/urls", d.GetUserURLs)

	r.Post("/", d.ShortenURL)
	r.Get("/{id}", d.ProcessShortURL)
	r.Get("/ping", d.Ping)

	return r
}
