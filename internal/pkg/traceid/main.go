// Package traceid Add/Get trace id from context.
package traceid

import "context"

type traceIDKey struct{}

// Add adds trace id to context.
func Add(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, id)
}

// Get gets trace id from context.
func Get(ctx context.Context) string {
	id, ok := ctx.Value(traceIDKey{}).(string)
	if !ok {
		return "undefined"
	}

	return id
}
