package session

import (
	"context"
)

// Define context keys to avoid collisions (use unexported type or package-private string)
type contextKey string

const (
	userIDKey  contextKey = "user_id"
	traceIDKey contextKey = "trace_id"
	// add more as needed
)

// ContextWithUserID Helper functions to store/retrieve
func ContextWithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}

func TraceIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(traceIDKey).(string)
	return id, ok
}
