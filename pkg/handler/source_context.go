package handler

import (
	"context"
)

type sourceContextKey struct{}

func contextWithSource(ctx context.Context, source string) context.Context {
	return context.WithValue(ctx, sourceContextKey{}, source)
}

func sourceFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(sourceContextKey{}).(string); ok {
		return v
	}
	return ""
}
