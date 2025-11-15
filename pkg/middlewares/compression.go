package middlewares

import (
	"compress/gzip"
	"errors"
	"net/http"
	"strings"
)

var (
	errGzipReaderInit = errors.New("failed to init a gzip reader")
)

var compressableContentTypes = []string{"application/json", "text/html"}

type gzipWriter struct {
	http.ResponseWriter
	gz *gzip.Writer
}

func (w gzipWriter) isCompressable() bool {
	contentTypeHeader := w.Header().Get("Content-Type")
	for _, ct := range compressableContentTypes {
		if strings.Contains(contentTypeHeader, ct) {
			return true
		}
	}
	return false
}

func (w gzipWriter) Write(b []byte) (int, error) {
	if w.isCompressable() {
		return w.gz.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w gzipWriter) WriteHeader(statusCode int) {
	if w.isCompressable() {
		w.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w gzipWriter) Close() error {
	if w.isCompressable() {
		return w.gz.Close()
	}
	return nil
}

func NewGzipCompressionMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				gz := gzipWriter{
					ResponseWriter: w,
					gz:             gzip.NewWriter(w),
				}
				w = gz
				defer gz.Close()
			}

			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				reader, err := gzip.NewReader(r.Body)
				if err != nil {
					e := errors.Join(err, errGzipReaderInit)
					http.Error(w, e.Error(), http.StatusInternalServerError)
					return
				}
				r.Body = reader
				defer reader.Close()
			}

			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
