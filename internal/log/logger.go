package log

import (
	"context"
	"github.com/IWannaWish/ethusd-converter/internal/config"
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...Field)
	Debug(ctx context.Context, msg string, args ...Field)
	Error(ctx context.Context, msg string, args ...Field)
	With(args ...Field) Logger
}

func NewLogger(cfg *config.Config) Logger {
	if cfg.UseZap {
		return NewZapLogger(cfg)
	}
	return NewSlogLogger(cfg)

}

type Field struct {
	Key   string
	Value any
}

func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Any(key string, value any) Field {
	return Field{Key: key, Value: value}
}
