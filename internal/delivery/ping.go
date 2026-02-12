package delivery

import "net/http"

// Пинг ручка для проверки работоспособности сервиса.
func (d *delivery) Ping(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()

	err := d.usecase.Ping(ctx)
	if err != nil {
		d.logger.WithError(err).Error("ping failed")
		http.Error(w, "ping failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
