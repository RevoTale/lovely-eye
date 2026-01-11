package auth

import (
	"context"
	"net/http"
)

// Service defines the authentication interface.
type Service interface {
	// Register creates a new user account. First user becomes admin.
	Register(ctx context.Context, input RegisterInput) (*User, *Tokens, error)

	// Login authenticates a user and returns tokens.
	Login(ctx context.Context, input LoginInput) (*User, *Tokens, error)

	// RefreshTokens generates new tokens from a valid refresh token.
	RefreshTokens(ctx context.Context, refreshToken string) (*Tokens, error)

	// ValidateAccessToken validates a token and returns claims.
	ValidateAccessToken(token string) (*Claims, error)

	// GetUserByID retrieves a user by ID.
	GetUserByID(ctx context.Context, id int64) (*User, error)

	// CreateInitialAdmin creates the first admin if configured.
	CreateInitialAdmin(ctx context.Context, username, password string) error

	// SetAuthCookies sets HTTP-only authentication cookies.
	SetAuthCookies(w http.ResponseWriter, tokens *Tokens)

	// ClearAuthCookies removes authentication cookies.
	ClearAuthCookies(w http.ResponseWriter)
}

// RegisterInput contains registration data.
type RegisterInput struct {
	Username string
	Password string
}

// LoginInput contains login credentials.
type LoginInput struct {
	Username string
	Password string
}

// Tokens contains authentication tokens.
type Tokens struct {
	AccessToken  string
	RefreshToken string
}

// User represents an authenticated user.
type User struct {
	ID       int64
	Username string
	Role     string
}

// Claims contains decoded token claims.
type Claims struct {
	UserID   int64
	Username string
	Role     string
}
