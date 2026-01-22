package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/delivery"
	"github.com/georgg2003/shortener/internal/repository/postgres"
	"github.com/georgg2003/shortener/internal/repository/storage"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	conf := config.New()
	if err := conf.ReadFromEnv(); err != nil {
		logger.WithError(err).Fatal("failed to read config from env")
	}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	if err := conf.ReadFromFlags(fs); err != nil {
		logger.WithError(err).Fatal("failed to read config from flags")
	}
	fs.Parse(os.Args[1:])

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var repo usecase.Repository
	if conf.DataBaseDSN != "" {
		repo = postgres.New(ctx, conf, logger)
	} else {
		repo = storage.New(ctx, conf, logger)
	}

	usecase := usecase.New(repo, conf, logger)
	delivery := delivery.New(usecase, logger, conf)

	r := delivery.GetNewRouter()

	log.Printf("Listening on %v", conf.ListenAddr)
	log.Fatal(http.ListenAndServe(conf.ListenAddr, r))
}
