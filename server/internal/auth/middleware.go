package auth

import (
	"context"
	"errors"
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
	jwtSvc, ok := service.(*jwtService)
	if !ok {
		panic("service must be a *jwtService")
	}
	return &Middleware{service: jwtSvc}
}

// Authenticate extracts and validates authentication from the request.
// Supports both cookie-based and header-based authentication.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var accessToken string

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
			if errors.Is(err, ErrExpiredToken) && refresh != "" {
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

func GetUserFromContext(ctx context.Context) *Claims {
	claims, ok := ctx.Value(UserContextKey).(*Claims)
	if !ok {
		return nil
	}
	return claims
}
