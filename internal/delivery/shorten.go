package delivery

import (
	"io"
	"net/http"
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

	shortUrl := d.usecase.NewShortUrl(string(bytes))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortUrl))
}
