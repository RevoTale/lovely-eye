package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("token has expired")
	ErrWrongTokenType = errors.New("wrong token type")
)

type tokenType string

const (
	accessTokenType  tokenType = "access"
	refreshTokenType tokenType = "refresh"
)

// jwtClaims is the internal JWT claims structure.
type jwtClaims struct {
	UserID    int64     `json:"uid"`
	Username  string    `json:"usr"`
	Role      string    `json:"rol"`
	TokenType tokenType `json:"typ"`
	jwt.RegisteredClaims
}

// jwtProvider handles JWT token generation and validation.
type jwtProvider struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	issuer        string
}

func newJWTProvider(secret string, accessExpiry, refreshExpiry time.Duration) *jwtProvider {
	return &jwtProvider{
		secret:        []byte(secret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		issuer:        "lovely-eye",
	}
}

func (p *jwtProvider) generateAccessToken(user *User) (string, error) {
	claims := &jwtClaims{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      user.Role,
		TokenType: accessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    p.issuer,
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}

func (p *jwtProvider) generateRefreshToken(user *User) (string, error) {
	claims := &jwtClaims{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      user.Role,
		TokenType: refreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    p.issuer,
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}

func (p *jwtProvider) validateToken(tokenString string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return p.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (p *jwtProvider) validateAccessToken(tokenString string) (*Claims, error) {
	jwtClaims, err := p.validateToken(tokenString)
	if err != nil {
		return nil, err
	}
	if jwtClaims.TokenType != accessTokenType {
		return nil, ErrWrongTokenType
	}
	return &Claims{
		UserID:   jwtClaims.UserID,
		Username: jwtClaims.Username,
		Role:     jwtClaims.Role,
	}, nil
}

func (p *jwtProvider) validateRefreshToken(tokenString string) (*Claims, error) {
	jwtClaims, err := p.validateToken(tokenString)
	if err != nil {
		return nil, err
	}
	if jwtClaims.TokenType != refreshTokenType {
		return nil, ErrWrongTokenType
	}
	return &Claims{
		UserID:   jwtClaims.UserID,
		Username: jwtClaims.Username,
		Role:     jwtClaims.Role,
	}, nil
}
