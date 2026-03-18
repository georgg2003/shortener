package interceptors

import (
	"context"

	"github.com/georgg2003/shortener/pkg/contextlib"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// Интерсептор для логирования в gRPC ручках.
func RequestIDInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	requestID := uuid.New().String()
	newCtx := contextlib.SetRequestID(ctx, requestID)
	return handler(newCtx, req)
}
