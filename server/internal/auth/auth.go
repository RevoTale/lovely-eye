package auth

import (
	"context"
	"net/http"
)

type Service interface {
	Register(ctx context.Context, input RegisterInput) (*User, *Tokens, error)

	Login(ctx context.Context, input LoginInput) (*User, *Tokens, error)

	RefreshTokens(ctx context.Context, refreshToken string) (*Tokens, error)

	ValidateAccessToken(token string) (*Claims, error)

	GetUserByID(ctx context.Context, id int64) (*User, error)

	CreateInitialAdmin(ctx context.Context, username, password string) error

	SetAuthCookies(w http.ResponseWriter, tokens *Tokens)

	ClearAuthCookies(w http.ResponseWriter)
}

type RegisterInput struct {
	Username string
	Password string
}

type LoginInput struct {
	Username string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
type User struct {
	ID       int64
	Username string
	Role     string
}

type Claims struct {
	UserID   int64
	Username string
	Role     string
}
