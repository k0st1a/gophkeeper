// Package userid Add/Get user id from context.
package userid

import "context"

type userIDKey struct{}

// Set adds user id to context.
func Set(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, userIDKey{}, id)
}

// Get gets user id from context.
func Get(ctx context.Context) (int64, bool) {
	i, ok := ctx.Value(userIDKey{}).(int64)
	return i, ok
}
