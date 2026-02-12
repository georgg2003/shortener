package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/georgg2003/shortener/internal/models"
)

// Ручка, возвращающая все URL, сокращенные пользователем.
func (d *delivery) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	defer d.safeClose(r.Body)

	ctx := r.Context()

	entities, err := d.usecase.GetUserURLs(ctx)
	if err != nil {
		d.logger.WithError(err).Error("get user urls failed")
		http.Error(w, "get user urls failed", http.StatusInternalServerError)
		return
	}

	if len(entities) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make(models.APIUserURLsResponse, 0, len(entities))
	for _, entity := range entities {
		resp = append(resp, models.APIUserURLsResponseRecord{
			OriginalURL: entity.OriginalURL,
			ShortURL:    entity.ShortURL,
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
