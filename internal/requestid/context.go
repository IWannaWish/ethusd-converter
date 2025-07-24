package requestid

import "context"

type contextKey string

const requestIDKey contextKey = "request_id"

func WithReqID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func ReqIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}
