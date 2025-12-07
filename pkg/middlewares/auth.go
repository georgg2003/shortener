package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/georgg2003/shortener/pkg/contextlib"
	"github.com/georgg2003/shortener/pkg/jwthelper"
	"github.com/sirupsen/logrus"
)

const tokenCookieName = "access_token"

type UseCase interface {
	NewUser(ctx context.Context) (userID int64, err error)
}

func createNewUser(ctx context.Context, uc UseCase, w http.ResponseWriter) (int64, error) {
	userID, err := uc.NewUser(ctx)
	if err != nil {
		return 0, err
	}
	token, err := jwthelper.NewAccessToken(userID)
	if err != nil {
		return 0, err
	}
	newCookie := &http.Cookie{
		Name:  tokenCookieName,
		Value: token,
	}
	http.SetCookie(w, newCookie)
	return userID, nil
}

func NewSimpleAuthMiddleware(l *logrus.Logger, uc UseCase) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var userID int64

			token, err := r.Cookie(tokenCookieName)
			if errors.Is(err, http.ErrNoCookie) {
				if userID, err = createNewUser(ctx, uc, w); err != nil {
					l.WithContext(ctx).WithError(err).Error("failed to create new user in middleware")
					http.Error(w, "failed to create a new user", http.StatusInternalServerError)
					return
				}
			} else {
				userID, err = jwthelper.ReadAccessToken(token.Value)
				if err != nil {
					l.WithContext(ctx).WithError(err).Error("got an invalid token")
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}
			}

			newCtx := contextlib.SetUserID(ctx, userID)

			h.ServeHTTP(w, r.WithContext(newCtx))
		}
		return http.HandlerFunc(fn)
	}
}
