package main

import (
	"log"
	"net/http"

	"github.com/georgg2003/shortener/internal/delivery"
	"github.com/georgg2003/shortener/internal/repository"
	"github.com/georgg2003/shortener/internal/usecase"
)

func main() {
	repo := repository.New()
	usecase := usecase.New(repo)
	delivery := delivery.New(usecase)

	r := delivery.GetNewRouter()

	log.Fatal(http.ListenAndServe(":8080", r))
}
