package log

import (
	"errors"
)

func WithStack(err error) []Field {
	if err == nil {
		return nil
	}
	fields := []Field{
		String("error", err.Error()),
	}
	if unwrapped := errors.Unwrap(err); unwrapped != nil {
		fields = append(fields, String("cause", unwrapped.Error()))
	}
	return fields
}
