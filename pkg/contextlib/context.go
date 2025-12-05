package contextlib

import "context"

var userKey = struct{}{}

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userKey, userID)
}

func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userKey).(int64)
	return userID, ok
}
