package applog

import (
	"context"
	"github.com/IWannaWish/ethusd-converter/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(cfg *config.Config) *ZapLogger {
	cfgZap := zap.NewDevelopmentConfig()
	cfgZap.Level = parseZapLevel(cfg.LogLevel)

	// Настройки вывода
	cfgZap.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfgZap.EncoderConfig.TimeKey = "time"
	cfgZap.OutputPaths = []string{"stdout"}

	logger, err := cfgZap.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		panic(err)
	}

	return &ZapLogger{logger: logger}
}

func (z *ZapLogger) Debug(ctx context.Context, msg string, args ...Field) {
	z.logger.Debug(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) Info(ctx context.Context, msg string, args ...Field) {
	z.logger.Info(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) Error(ctx context.Context, msg string, args ...Field) {
	z.logger.Error(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) With(args ...Field) Logger {
	return &ZapLogger{
		logger: z.logger.With(z.toZapFields(context.Background(), args)...),
	}
}

func (z *ZapLogger) toZapFields(ctx context.Context, fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields)+1)

	if reqID, ok := RequestIDFromContext(ctx); ok {
		zapFields = append(zapFields, zap.String("request_id", reqID))
	}

	for _, f := range fields {
		switch v := f.Value.(type) {
		case string:
			zapFields = append(zapFields, zap.String(f.Key, v))
		case error:
			zapFields = append(zapFields, zap.NamedError(f.Key, v))
		default:
			zapFields = append(zapFields, zap.Any(f.Key, v))
		}
	}

	return zapFields
}

func parseZapLevel(level string) zap.AtomicLevel {
	switch level {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}
