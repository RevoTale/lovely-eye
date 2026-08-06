package auth

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

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
	require.False(t, first.CreatedAt.IsZero())

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

func TestAuthServiceDoesNotClassifyStorageFailureAsUserError(t *testing.T) {
	storageErr := errors.New("storage unavailable")
	service := NewService(&fakeUserStore{err: storageErr}, testAuthConfig(true))

	_, _, loginErr := service.Login(context.Background(), LoginInput{
		Username: "missing",
		Password: "password123",
	})
	require.ErrorIs(t, loginErr, storageErr)
	require.False(t, errors.Is(loginErr, ErrInvalidCredentials))

	_, userErr := service.GetUserByID(context.Background(), 999)
	require.ErrorIs(t, userErr, storageErr)
	require.False(t, errors.Is(userErr, ErrUserNotFound))
}

func newTestAuthService(t *testing.T, allowRegistration bool) *Service {
	t.Helper()
	return NewService(newFakeUserStore(), testAuthConfig(allowRegistration))
}

func testAuthConfig(allowRegistration bool) Config {
	return Config{
		JWTSecret:         strings.Repeat("j", 32),
		AccessTokenExpiry: 15 * time.Minute,
		RefreshExpiry:     7 * 24 * time.Hour,
		AllowRegistration: allowRegistration,
	}
}

type fakeUserStore struct {
	users  map[int64]*StoredUser
	nextID int64
	err    error
}

func newFakeUserStore() *fakeUserStore {
	return &fakeUserStore{users: make(map[int64]*StoredUser), nextID: 1}
}

func (s *fakeUserStore) CreateForRegistration(
	_ context.Context,
	user *StoredUser,
	allowRegistration bool,
) error {
	if s.err != nil {
		return s.err
	}
	if len(s.users) != 0 && !allowRegistration {
		return ErrRegistrationDisabled
	}
	for _, existing := range s.users {
		if existing.Username == user.Username {
			return ErrUserExists
		}
	}
	user.ID = s.nextID
	user.Role = "user"
	if len(s.users) == 0 {
		user.Role = "admin"
	}
	user.CreatedAt = time.Now()
	s.nextID++
	s.users[user.ID] = cloneStoredUser(user)
	return nil
}

func (s *fakeUserStore) CreateInitialAdminIfNoUsers(
	_ context.Context,
	user *StoredUser,
) (bool, error) {
	if s.err != nil {
		return false, s.err
	}
	if len(s.users) != 0 {
		return false, nil
	}
	user.ID = s.nextID
	user.Role = "admin"
	user.CreatedAt = time.Now()
	s.nextID++
	s.users[user.ID] = cloneStoredUser(user)
	return true, nil
}

func (s *fakeUserStore) GetByID(_ context.Context, id int64) (*StoredUser, error) {
	if s.err != nil {
		return nil, s.err
	}
	user, ok := s.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return cloneStoredUser(user), nil
}

func (s *fakeUserStore) GetByUsername(_ context.Context, username string) (*StoredUser, error) {
	if s.err != nil {
		return nil, s.err
	}
	for _, user := range s.users {
		if user.Username == username {
			return cloneStoredUser(user), nil
		}
	}
	return nil, ErrUserNotFound
}

func (s *fakeUserStore) HasUsers(context.Context) (bool, error) {
	if s.err != nil {
		return false, s.err
	}
	return len(s.users) != 0, nil
}

func cloneStoredUser(user *StoredUser) *StoredUser {
	clone := *user
	return &clone
}
