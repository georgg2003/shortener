package delivery

import (
	"net/http"

	"github.com/georgg2003/shortener/internal/usecase"
)

type Delivery interface {
	ShortenURL(w http.ResponseWriter, r *http.Request)
	ProcessShortURL(w http.ResponseWriter, r *http.Request)
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
