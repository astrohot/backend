package auth

import (
	"context"
)

type (
	// A private key for context that only this package can access. This is important
	// to prevent collisions between different context uses
	contextKey struct{}

	// Auther is an interface for entities that need to be authenticated. It
	// receives a plain string to be compared with the hash password. It
	// returns a valid token string or an error.
	Auther interface {
		Authenticate(string) (string, error)
	}
)

// WithContext ...
func WithContext(ctx context.Context, entity Auther) context.Context {
	return context.WithValue(ctx, contextKey{}, entity)
}

// FromContext ...
func FromContext(ctx context.Context) Auther {
	if a, ok := ctx.Value(contextKey{}).(Auther); ok {
		return a
	}

	return nil
}
