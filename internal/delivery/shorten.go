package delivery

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (d delivery) ShortenURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bytes, err := io.ReadAll(r.Body)
	if err != nil || len(bytes) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	longURL := string(bytes)

	parsedURL, parseErr := url.Parse(longURL)

	if parseErr != nil {
		http.Error(w, "Invalid url", http.StatusBadRequest)
		return
	}
	if parsedURL.Scheme == "" {
		http.Error(w, "Invalid scheme", http.StatusBadRequest)
		return
	}
	if parsedURL.Hostname() == "" {
		http.Error(w, "Invalid hostname", http.StatusBadRequest)
		return
	}

	shortURL := d.usecase.NewShortURL(longURL)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(shortURL)))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}
