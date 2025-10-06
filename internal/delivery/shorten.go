package delivery

import (
	"io"
	"net/http"
	"net/url"
)

func (d delivery) ShortenURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil || len(bytes) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	longUrl := string(bytes)

	parsedUrl, parseErr := url.Parse(longUrl)

	if parseErr != nil {
		http.Error(w, "Invalid url", http.StatusBadRequest)
		return
	}
	if parsedUrl.Scheme == "" {
		http.Error(w, "Invalid scheme", http.StatusBadRequest)
		return
	}
	if parsedUrl.Hostname() == "" {
		http.Error(w, "Invalid hostname", http.StatusBadRequest)
		return
	}

	shortUrl := d.usecase.NewShortUrl(longUrl)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortUrl))
}
