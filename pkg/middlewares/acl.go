package middlewares

import (
	"net"
	"net/http"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/sirupsen/logrus"
)

const RealIPHeaderName = "X-Real-IP"

// Проверяет IP адрес пришедшего запроса. Возвращает 403, если запрос пришел из недоверенной сети.
func NewACLMiddleware(
	cfg *config.Config,
	l logrus.FieldLogger,
) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			subnet := cfg.TrustedSubnet
			realIP := net.ParseIP(r.Header.Get(RealIPHeaderName))
			if subnet == nil || realIP == nil || !subnet.Contains(realIP) {
				l.WithFields(logrus.Fields{
					"trusted_subnet": subnet,
					"real_ip":        realIP,
				}).Error("access denied for internal api")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
