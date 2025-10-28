package main

import (
	"log"
	"net/http"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/delivery"
	"github.com/georgg2003/shortener/internal/repository"
	"github.com/georgg2003/shortener/internal/usecase"
)

func main() {
	conf := config.New()
	conf.ReadFromEnv()
	conf.ReadFromFlags()

	repo := repository.New()
	usecase := usecase.New(repo, conf)
	delivery := delivery.New(usecase)

	r := delivery.GetNewRouter()

	log.Printf("Listening on %v", conf.ListenAddr)
	log.Fatal(http.ListenAndServe(conf.ListenAddr, r))
}
