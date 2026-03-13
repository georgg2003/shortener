package interceptors

import (
	"context"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/georgg2003/shortener/pkg/contextlib"
	"github.com/georgg2003/shortener/pkg/jwthelper"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const AuthKeyName = "authorization"

func createNewUser(
	ctx context.Context,
	uc usecase.UseCase,
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
	md := metadata.Pairs(AuthKeyName, token)
	grpc.SetHeader(ctx, md)
	return userID, nil
}

// Интерсептор для авторизации в gRPC ручках.
// Авторизует по JWT токену в куках.
// Если пользак не авторизован, создает нового.
// Текущего пользователя записывает в контекст.
func NewAuthInterceptor(cfg *config.Config, l logrus.FieldLogger, uc usecase.UseCase) grpc.UnaryServerInterceptor {
	helper := jwthelper.New([]byte(cfg.JWTSecretKey))

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		var userID int64
		values := metadata.ValueFromIncomingContext(ctx, AuthKeyName)
		if len(values) > 0 {
			token := values[0]
			userID, err = helper.ReadAccessToken(token)
			if err != nil {
				l.WithError(err).Error("got an invalid token")
				return nil, status.Error(codes.Unauthenticated, "invalid token")
			}
		} else {
			if userID, err = createNewUser(ctx, uc, helper); err != nil {
				l.WithError(err).Error("failed to create new user in middleware")
				return
			}
		}

		newCtx := contextlib.SetUserID(ctx, userID)
		return handler(newCtx, req)
	}
}
