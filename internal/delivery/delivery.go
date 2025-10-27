package delivery

import (
	"net/http"

	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type Delivery interface {
	ShortenURL(w http.ResponseWriter, r *http.Request)
	ProcessShortURL(w http.ResponseWriter, r *http.Request)

	GetNewRouter() chi.Router
}

type delivery struct {
	usecase usecase.UseCase
}

func New(
	usecase usecase.UseCase,
) Delivery {
	return delivery{
		usecase: usecase,
	}
}

func (d delivery) GetNewRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", d.ShortenURL)
	r.Get("/{id}", d.ProcessShortURL)

	return r
}
