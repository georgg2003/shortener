package delivery

import (
	"errors"
	"net/http"

	"github.com/georgg2003/shortener/internal/repository/db"
)

func (d *delivery) ProcessShortURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	longURL, isDeleted, err := d.usecase.ProcessShortURL(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "Link not found", http.StatusNotFound)
			return
		}
		d.logger.WithContext(ctx).WithError(err).Error("failed to process short url")
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
