package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTProviderRejectsUnexpectedSigningMethod(t *testing.T) {
	provider := newJWTProvider(strings.Repeat("s", 32), 15*time.Minute, 7*24*time.Hour)
	claims := validTestClaims(provider, accessTokenType)
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	signed, err := token.SignedString(provider.secret)
	require.NoError(t, err)

	_, err = provider.validateAccessToken(signed)
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestJWTProviderRejectsUnexpectedIssuer(t *testing.T) {
	provider := newJWTProvider(strings.Repeat("s", 32), 15*time.Minute, 7*24*time.Hour)
	claims := validTestClaims(provider, accessTokenType)
	claims.Issuer = "other-issuer"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(provider.secret)
	require.NoError(t, err)

	_, err = provider.validateAccessToken(signed)
	require.ErrorIs(t, err, ErrInvalidToken)
}

func validTestClaims(provider *jwtProvider, tokenType tokenType) *jwtClaims {
	now := time.Now()
	return &jwtClaims{
		UserID:    1,
		Username:  "admin",
		Role:      "admin",
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(provider.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    provider.issuer,
			Subject:   "admin",
		},
	}
}
