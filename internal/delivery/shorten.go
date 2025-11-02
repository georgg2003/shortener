package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/georgg2003/shortener/internal/models"
)

var (
	invalidURLError         = errors.New("Invalid url")
	invalidURLSchemeError   = errors.New("Invalid scheme")
	invalidURLHostnameError = errors.New("Invalid hostname")

	decodeBodyError = errors.New("failed to decode body")
)

func validateURL(longURL string) error {
	parsedURL, parseErr := url.Parse(longURL)

	if parseErr != nil {
		return errors.Join(parseErr, invalidURLError)
	}
	if parsedURL.Scheme == "" {
		return invalidURLSchemeError
	}
	if parsedURL.Hostname() == "" {
		return invalidURLHostnameError
	}

	return nil
}

func (d delivery) ShortenURL(w http.ResponseWriter, r *http.Request) {
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
	}

	shortURL := d.usecase.NewShortURL(ctx, longURL)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(shortURL)))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (d delivery) APIShortenURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	var req models.APIShortenURLRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		e := errors.Join(err, decodeBodyError)
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	if err := validateURL(req.URL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := d.usecase.NewShortURL(ctx, req.URL)

	response := models.APIShortenURLResponse{
		Result: shortURL,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
