package delivery

import (
	"errors"
	"net/http"

	"github.com/georgg2003/shortener/internal/repository/db"
)

// Ручка переадресации на оригинальный URL по сокращенному.
func (d *delivery) ProcessShortURL(w http.ResponseWriter, r *http.Request) {
	defer d.safeClose(r.Body)
	ctx := r.Context()

	id := r.PathValue("id")
	longURL, isDeleted, err := d.usecase.ProcessShortURL(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "Link not found", http.StatusNotFound)
			return
		}
		d.logger.WithError(err).Error("failed to process short url")
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	if isDeleted {
		http.Error(w, "URL was deleted", http.StatusGone)
		return
	}

	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
