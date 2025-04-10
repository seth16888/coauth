package middleware

import (
	"coauth/pkg/helpers"
	"context"

	"google.golang.org/grpc"
)

func RequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		// 生成Request ID
		requestID := helpers.NewUUID()

		ctx = context.WithValue(ctx, "X-Request-Id", requestID)

		// 继续处理请求
		return handler(ctx, req)
	}
}
