package auth

import "context"

type contextKey struct{}

// ContextWithClaims attaches validated authentication claims to a request context.
func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, contextKey{}, claims)
}

// GetUserFromContext returns validated claims when the request is authenticated.
func GetUserFromContext(ctx context.Context) *Claims {
	claims, ok := ctx.Value(contextKey{}).(*Claims)
	if !ok {
		return nil
	}
	return claims
}
