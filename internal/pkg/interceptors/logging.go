package interceptors

import (
	"context"
	"time"

	"github.com/georgg2003/shortener/pkg/contextlib"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Интерсептор для логирования в gRPC ручках.
func NewLoggingInterceptor(l logrus.FieldLogger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		now := time.Now()
		reqID, _ := contextlib.GetRequestID(ctx)

		l.WithFields(logrus.Fields{
			"method":     info.FullMethod,
			"request_id": reqID,
		}).Info("new grpc request")

		resp, err = handler(ctx, req)
		duration := time.Since(now)

		l.WithFields(logrus.Fields{
			"method":     info.FullMethod,
			"duraton":    duration.String(),
			"error":      err,
			"request_id": reqID,
		}).Info("grpc response")

		return resp, err
	}
}
