package auth

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/lovely-eye/server/internal/config"
	"github.com/lovely-eye/server/internal/database"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestRegisterValidatesCredentialsAtServerBoundary(t *testing.T) {
	service := newTestAuthService(t, true)

	tests := []struct {
		name  string
		input RegisterInput
		err   error
	}{
		{
			name:  "empty username",
			input: RegisterInput{Username: "   ", Password: "password123"},
			err:   ErrInvalidUsername,
		},
		{
			name:  "control character username",
			input: RegisterInput{Username: "bad\nname", Password: "password123"},
			err:   ErrInvalidUsername,
		},
		{
			name:  "short password",
			input: RegisterInput{Username: "admin", Password: "short"},
			err:   ErrInvalidPassword,
		},
		{
			name:  "bcrypt truncation boundary",
			input: RegisterInput{Username: "admin", Password: strings.Repeat("p", 73)},
			err:   ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.Register(context.Background(), tt.input)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestRegisterNormalizesUsernameAndAssignsOnlyFirstUserAdmin(t *testing.T) {
	service := newTestAuthService(t, true)
	ctx := context.Background()

	first, _, err := service.Register(ctx, RegisterInput{
		Username: "  admin  ",
		Password: "password123",
	})
	require.NoError(t, err)
	require.Equal(t, "admin", first.Username)
	require.Equal(t, "admin", first.Role)

	second, _, err := service.Register(ctx, RegisterInput{
		Username: "user",
		Password: "password123",
	})
	require.NoError(t, err)
	require.Equal(t, "user", second.Role)

	_, _, err = service.Register(ctx, RegisterInput{
		Username: "admin",
		Password: "password123",
	})
	require.ErrorIs(t, err, ErrUserExists)
}

func TestRegisterRejectsSecondUserWhenRegistrationDisabled(t *testing.T) {
	service := newTestAuthService(t, false)
	ctx := context.Background()

	_, _, err := service.Register(ctx, RegisterInput{
		Username: "admin",
		Password: "password123",
	})
	require.NoError(t, err)

	_, _, err = service.Register(ctx, RegisterInput{
		Username: "user",
		Password: "password123",
	})
	require.ErrorIs(t, err, ErrRegistrationDisabled)

	_, _, err = service.Register(ctx, RegisterInput{
		Username: "admin",
		Password: "password123",
	})
	require.ErrorIs(t, err, ErrRegistrationDisabled)
}

func TestCreateInitialAdminValidatesCredentials(t *testing.T) {
	service := newTestAuthService(t, true)

	err := service.CreateInitialAdmin(context.Background(), "admin", "short")
	require.ErrorIs(t, err, ErrInvalidPassword)
}

func newTestAuthService(t *testing.T, allowRegistration bool) *jwtService {
	t.Helper()

	dbName := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	db, err := database.New(config.DatabaseConfig{
		Driver:         config.DBDriverSQLite,
		DSN:            fmt.Sprintf("file:%s?mode=memory&cache=shared", dbName),
		ConnectTimeout: time.Second,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, database.Close(db))
	})
	require.NoError(t, database.Migrate(context.Background(), db))

	repo := repository.NewUserRepository(db)
	return NewService(repo, Config{
		JWTSecret:         strings.Repeat("j", 32),
		AccessTokenExpiry: 15 * time.Minute,
		RefreshExpiry:     7 * 24 * time.Hour,
		AllowRegistration: allowRegistration,
		SecureCookies:     false,
	}).(*jwtService)
}
