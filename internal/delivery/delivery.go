package delivery

import (
	"net/http"

	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/georgg2003/shortener/pkg/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Delivery interface {
	ShortenURL(w http.ResponseWriter, r *http.Request)
	ProcessShortURL(w http.ResponseWriter, r *http.Request)

	GetNewRouter() chi.Router
}

type delivery struct {
	usecase usecase.UseCase
	logger  *logrus.Logger
}

func New(
	usecase usecase.UseCase,
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
