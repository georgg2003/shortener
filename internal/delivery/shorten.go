package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/internal/usecase"
)

var (
	errInvalidURL        = errors.New("invalid url")
	errInvalidURLScheme  = errors.New("invalid scheme")
	errIvalidURLHostname = errors.New("invalid hostname")

	errDecodeBody = errors.New("failed to decode body")
)

func validateURL(longURL string) error {
	parsedURL, parseErr := url.Parse(longURL)

	if parseErr != nil {
		return errors.Join(parseErr, errInvalidURL)
	}
	if parsedURL.Scheme == "" {
		return errInvalidURLScheme
	}
	if parsedURL.Hostname() == "" {
		return errIvalidURLHostname
	}

	return nil
}

// Ручка сокращения URL.
func (d *delivery) ShortenURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	bytes, err := io.ReadAll(r.Body)
	if err != nil || len(bytes) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	longURL := string(bytes)

	if err = validateURL(longURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := d.usecase.NewShortURL(ctx, longURL)
	if err != nil && !errors.Is(err, usecase.ErrURLEntityAlreadyExists) {
		d.logger.WithError(err).Error("failed to add a new short url")
		http.Error(w, internalErrorText, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(shortURL)))

	if errors.Is(err, usecase.ErrURLEntityAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(shortURL))
}

// API ручка для сокращения URL.
func (d *delivery) APIShortenURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	var req models.APIShortenURLRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		e := errors.Join(err, errDecodeBody)
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	if err := validateURL(req.URL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := d.usecase.NewShortURL(ctx, req.URL)
	if err != nil && !errors.Is(err, usecase.ErrURLEntityAlreadyExists) {
		d.logger.WithError(err).Error("failed to create a new short url")
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	response := models.APIShortenURLResponse{
		Result: shortURL,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if errors.Is(err, usecase.ErrURLEntityAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// API ручка для сокращения нескольких URL.
func (d *delivery) APIShortenURLBatch(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	var req models.APIShortenURLBatchRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		e := errors.Join(err, errDecodeBody)
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	entities := make([]*models.URLEntity, 0, len(req))
	for _, v := range req {
		if err := validateURL(v.OriginalURL); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		entities = append(entities, &models.URLEntity{
			CorrelationID: v.CorrelationID,
			OriginalURL:   v.OriginalURL,
		})
	}

	err := d.usecase.NewShortURLBatch(ctx, entities)
	if err != nil && !errors.Is(err, usecase.ErrURLEntityAlreadyExists) {
		d.logger.WithError(err).Error("failed to add batch of urls")
		http.Error(w, internalErrorText, http.StatusInternalServerError)
		return
	}

	response := make(models.APIShortenURLBatchResponse, 0, len(entities))
	for _, v := range entities {
		response = append(response, models.APIShortenURLBatchResponseRecord{
			CorrelationID: v.CorrelationID,
			ShortURL:      v.ShortURL,
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if errors.Is(err, usecase.ErrURLEntityAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
