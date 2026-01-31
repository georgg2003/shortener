package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/pkg/contextlib"
	"github.com/georgg2003/shortener/pkg/jwthelper"
	"github.com/sirupsen/logrus"
)

const tokenCookieName = "access_token"

type UseCase interface {
	NewUser(ctx context.Context) (userID int64, err error)
}

func createNewUser(
	ctx context.Context,
	uc UseCase,
	w http.ResponseWriter,
	helper *jwthelper.JWTHelper,
) (int64, error) {
	userID, err := uc.NewUser(ctx)
	if err != nil {
		return 0, err
	}
	token, err := helper.NewAccessToken(userID)
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

// Авторизационная мидлваря для http сервера.
// Авторизует по JWT токену в куках.
// Если пользак не авторизован, создает нового.
// Текущего пользователя записывает в контекст.
func NewSimpleAuthMiddleware(
	cfg *config.Config,
	l *logrus.Logger,
	uc UseCase,
) func(h http.Handler) http.Handler {
	helper := jwthelper.New([]byte(cfg.JWTSecretKey))
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var userID int64

			token, err := r.Cookie(tokenCookieName)
			if errors.Is(err, http.ErrNoCookie) {
				if userID, err = createNewUser(ctx, uc, w, helper); err != nil {
					l.WithContext(ctx).WithError(err).Error("failed to create new user in middleware")
					http.Error(w, "failed to create a new user", http.StatusInternalServerError)
					return
				}
			} else {
				userID, err = helper.ReadAccessToken(token.Value)
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
