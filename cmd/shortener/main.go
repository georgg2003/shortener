package main

import (
	"net/http"

	"github.com/georgg2003/shortener/internal/delivery"
	"github.com/georgg2003/shortener/internal/repository"
	"github.com/georgg2003/shortener/internal/usecase"
)

func main() {
	repo := repository.New()
	usecase := usecase.New(repo)
	delivery := delivery.New(usecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/", delivery.ShortenURL)
	mux.HandleFunc("/{id}/", delivery.ProcessShortURL)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
