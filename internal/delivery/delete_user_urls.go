package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/utils"
)

// Ручка удаления нескольких сокращенных URL пользователя.
// Принимает список сокращенных ID.
func (d *delivery) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()

	var req models.APIDeleteUserURLsRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		e := utils.ErrWrap(err, errDecodeBody.Error())
		d.logger.WithError(e).Info("bad request")
		http.Error(w, "failed to decode body", http.StatusBadRequest)
		return
	}

	if err := d.usecase.DeleteUserURLs(ctx, req); err != nil {
		d.logger.WithError(err).Error("delete user urls internal err")
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
