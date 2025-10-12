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
	conf, err := config.ReadFromYaml()
	if err != nil || conf == nil {
		log.Fatal("Failed to init config")
	}

	repo := repository.New()
	usecase := usecase.New(repo, *conf)
	delivery := delivery.New(usecase)

	r := delivery.GetNewRouter()

	log.Fatal(http.ListenAndServe(conf.ListenAddr, r))
}
