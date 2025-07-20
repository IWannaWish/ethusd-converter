package log

import (
	"context"
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger() *SlogLogger {
	format := os.Getenv("LOG_FORMAT") // json/text
	level := parseLogLevel(os.Getenv("LOG_LEVEL"))

	var handler slog.Handler

	if format == "text" {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level: slog.LevelDebug,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}
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
	if reqID, ok := RequestIDFromContext(ctx); ok {
		args = append(args, "request_id", reqID)
	}
	for _, f := range fields {
		args = append(args, f.Key, f.Value)
	}
	return args
}
