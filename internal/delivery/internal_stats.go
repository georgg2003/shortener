package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/georgg2003/shortener/internal/models"
)

func (d delivery) InternalStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, err := d.usecase.GetStats(ctx)
	if err != nil {
		d.logger.WithError(err).Error("failed to get stats")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := models.APIInternalStatsResponse{
		URLsCount:  stats.URLsCount,
		UsersCount: stats.UsersCount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		d.logger.WithError(err).Error("failed to write stats response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
