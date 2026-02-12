package middlewares

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type responseMetaData struct {
	statusCode int
	bodySize   int
}

type extendedResponseWriter struct {
	http.ResponseWriter
	responseMetaData *responseMetaData
}

func (w extendedResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.responseMetaData.bodySize += size
	return size, err
}

func (w extendedResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.responseMetaData.statusCode = statusCode
}

// Мидлваря, которая записывает в логи запросы и ответы.
func NewAccessLogMiddleware(l *logrus.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			method := r.Method
			uri := r.RequestURI

			l.WithFields(logrus.Fields{
				"method": method,
				"uri":    uri,
			}).Info("Request")

			responseMetaData := responseMetaData{}

			writer := extendedResponseWriter{
				ResponseWriter:   w,
				responseMetaData: &responseMetaData,
			}

			h.ServeHTTP(writer, r)

			duration := time.Since(now)

			l.WithFields(logrus.Fields{
				"duration":    duration.String(),
				"status_code": responseMetaData.statusCode,
				"body_size":   responseMetaData.bodySize,
			}).Info("Response")
		}
		return http.HandlerFunc(fn)
	}
}
