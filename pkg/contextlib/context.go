// Пакет для работы с контекстом.
package contextlib

import "context"

type userKey struct{}

var uk = userKey{}

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, uk, userID)
}

func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(uk).(int64)
	return userID, ok
}

type requestIDKey struct{}

var reqIDKey = requestIDKey{}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, reqIDKey, requestID)
}

func GetRequestID(ctx context.Context) (string, bool) {
	reqID, ok := ctx.Value(reqIDKey).(string)
	return reqID, ok
}
