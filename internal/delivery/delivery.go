package delivery

import (
	"context"

	"github.com/georgg2003/shortener/pkg/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Delivery interface {
	GetNewRouter() chi.Router
}

type UseCase interface {
	NewShortURL(ctx context.Context, url string) string
	ProcessShortURL(ctx context.Context, id string) (string, error)
	Ping(ctx context.Context) error
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
		middlewares.NewGzipCompressionMiddleware(),
	)

	r.Post("/api/shorten", d.APIShortenURL)

	r.Post("/", d.ShortenURL)
	r.Get("/{id}", d.ProcessShortURL)
	r.Get("/ping", d.Ping)

	return r
}
