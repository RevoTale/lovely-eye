package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserContextKey contextKey = "auth_user"
)

// Middleware provides HTTP middleware for authentication.
type Middleware struct {
	service *jwtService
}

// NewMiddleware creates a new authentication middleware.
func NewMiddleware(service Service) *Middleware {
	return &Middleware{service: service.(*jwtService)}
}

// Authenticate extracts and validates authentication from the request.
// Supports both cookie-based and header-based authentication.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var accessToken string

		// Try cookie first
		access, refresh := m.service.getTokensFromRequest(r)
		accessToken = access

		// Fallback to Authorization header
		if accessToken == "" {
			if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
				accessToken = strings.TrimPrefix(h, "Bearer ")
			}
		}

		if accessToken == "" {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := m.service.ValidateAccessToken(accessToken)
		if err != nil {
			// Try refresh if access token expired
			if err == ErrExpiredToken && refresh != "" {
				if tokens, refreshErr := m.service.RefreshTokens(r.Context(), refresh); refreshErr == nil {
					m.service.SetAuthCookies(w, tokens)
					claims, _ = m.service.ValidateAccessToken(tokens.AccessToken)
				}
			}
			if claims == nil {
				next.ServeHTTP(w, r)
				return
			}
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireCSRF validates the CSRF token for state-changing requests.
func (m *Middleware) RequireCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip for safe methods
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// Skip for API clients using Authorization header
		if r.Header.Get("Authorization") != "" {
			next.ServeHTTP(w, r)
			return
		}

		// Validate CSRF for cookie-based auth
		if !m.service.validateCSRF(r) {
			http.Error(w, "Invalid CSRF token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext extracts user claims from request context.
func GetUserFromContext(ctx context.Context) *Claims {
	claims, ok := ctx.Value(UserContextKey).(*Claims)
	if !ok {
		return nil
	}
	return claims
}
