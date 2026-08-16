package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/lovely-eye/server/internal/auth"
)

type tokenService interface {
	ValidateAccessToken(token string) (*auth.Claims, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*auth.Tokens, error)
}

type authMiddleware struct {
	service tokenService
	cookies *CookieManager
}

func newAuthMiddleware(service tokenService, cookies *CookieManager) *authMiddleware {
	return &authMiddleware{service: service, cookies: cookies}
}

// authenticate extracts and validates HttpOnly cookie authentication.
func (m *authMiddleware) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		access, refresh := m.cookies.TokensFromRequest(r)
		if access == "" {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := m.service.ValidateAccessToken(access)
		if err != nil {
			if errors.Is(err, auth.ErrExpiredToken) && refresh != "" {
				if tokens, refreshErr := m.service.RefreshTokens(r.Context(), refresh); refreshErr == nil {
					m.cookies.SetAuthCookies(w, tokens)
					claims, _ = m.service.ValidateAccessToken(tokens.AccessToken)
				}
			}
			if claims == nil {
				next.ServeHTTP(w, r)
				return
			}
		}

		ctx := auth.ContextWithClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
