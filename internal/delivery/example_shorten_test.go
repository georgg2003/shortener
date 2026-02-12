package delivery_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/delivery"
	"github.com/georgg2003/shortener/internal/repository/db/storage"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/georgg2003/shortener/pkg/utils"
	"github.com/sirupsen/logrus"
)

func ExampleDelivery_ShortenURL() {
	conf := config.New()
	logger := logrus.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := storage.New(ctx, conf, logger)
	usecase := usecase.New(repo, conf, logger, usecase.WithBase62Generator(utils.StubGenerator{}))
	delivery := delivery.New(usecase, logger, conf)

	buf := bytes.NewBufferString("https://calendar.mail.ru")
	req, err := http.NewRequest("POST", "/", buf)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(delivery.ShortenURL)
	handler.ServeHTTP(rr, req)
	fmt.Println(rr.Code)
	fmt.Println(rr.Body)
	// OUTPUT:
	// 201
	// http://localhost:8080/01234567
}
