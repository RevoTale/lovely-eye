package auth

import (
	"context"
	"time"
)

// UserStore is the persistence contract required by authentication.
// Adapters own database details and translate storage failures into auth errors.
type UserStore interface {
	CreateForRegistration(ctx context.Context, user *StoredUser, allowRegistration bool) error
	CreateInitialAdminIfNoUsers(ctx context.Context, user *StoredUser) (bool, error)
	GetByID(ctx context.Context, id int64) (*StoredUser, error)
	GetByUsername(ctx context.Context, username string) (*StoredUser, error)
	HasUsers(ctx context.Context) (bool, error)
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
	ID        int64
	Username  string
	Role      string
	CreatedAt time.Time
}

// StoredUser contains the authentication data that persistence must retain.
// PasswordHash intentionally never crosses the public auth service boundary.
type StoredUser struct {
	ID           int64
	Username     string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
}

type RegistrationStatus struct {
	HasUsers          bool
	AllowRegistration bool
}

type Claims struct {
	UserID   int64
	Username string
	Role     string
}
