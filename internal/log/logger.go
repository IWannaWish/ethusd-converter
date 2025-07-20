package log

import "context"

type Logger interface {
	Info(ctx context.Context, msg string, args ...Field)
	Debug(ctx context.Context, msg string, args ...Field)
	Error(ctx context.Context, msg string, args ...Field)
	With(args ...Field) Logger
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
