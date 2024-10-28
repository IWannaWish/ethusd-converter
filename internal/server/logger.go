package server

import (
	"context"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func FetchLoggingInterceptor() func(context.Context, any, *grpc.UnaryServerInfo, grpc.UnaryHandler) (resp any, err error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		logger := log.With().Str("method", info.FullMethod).Logger()

		logger.Info().Msg("request started")

		hResult, err := handler(ctx, req)

		st, ok := status.FromError(err)
		if !ok {
			logger.Warn().Msg("handler returned unknown status code")
		}

		l := logger.Info().Stringer("code", st.Code())
		if st.Message() != "" {
			l = l.Str("descr", st.Message())
		}
		l.Msg("request completed")
		return hResult, err
	}
}
