package grpc_server

import (
	"context"
	"github.com/IWannaWish/ethusd-converter/internal/applog"
	"github.com/IWannaWish/ethusd-converter/internal/requestid"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const requestIDHeader = "x-request-id"

func RequestIDInterceptor(logger applog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		var rid string

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get(requestIDHeader); len(values) > 0 && values[0] != "" {
				rid = values[0]
			}
		}
		if rid == "" {
			rid = uuid.NewString()
		}
		ctx = requestid.WithRequestID(ctx, rid)
		logger.Debug(ctx, "requestID установлен через interceptor", applog.String("request_id", rid))

		return handler(ctx, req)
	}
}
