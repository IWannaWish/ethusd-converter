package applog

import (
	"context"
	"github.com/IWannaWish/ethusd-converter/internal/config"
	"github.com/IWannaWish/ethusd-converter/internal/requestid"
	"log/slog"
	"os"
	"strings"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger(cfg *config.Config) *SlogLogger {
	level := parseLogLevel(cfg.LogLevel)

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "stack" {
				return slog.Attr{}
			}
			if a.Key == "request_id" {
				a.Key = "rid"
			}
			return a
		},
	})

	return &SlogLogger{
		logger: slog.New(handler),
	}
}

func (s SlogLogger) Info(ctx context.Context, msg string, args ...Field) {
	s.logger.InfoContext(ctx, msg, toArgs(ctx, args)...)
}

func (s SlogLogger) Debug(ctx context.Context, msg string, args ...Field) {
	s.logger.DebugContext(ctx, msg, toArgs(ctx, args)...)
}

func (s SlogLogger) Error(ctx context.Context, msg string, args ...Field) {
	s.logger.ErrorContext(ctx, msg, toArgs(ctx, args)...)
}

func (s SlogLogger) With(args ...Field) Logger {
	attrArgs := make([]any, 0, len(args)*2)
	for _, f := range args {
		attrArgs = append(attrArgs, f.Key, f.Value)
	}
	return &SlogLogger{
		logger: s.logger.With(attrArgs...),
	}
}

func toArgs(ctx context.Context, fields []Field) []any {
	args := make([]any, 0, len(fields)*2)
	if reqID, ok := requestid.RequestIDFromContext(ctx); ok {
		args = append(args, "request_id", reqID)
	}
	for _, f := range fields {
		args = append(args, f.Key, f.Value)
	}
	return args
}

func parseLogLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
