package persistence

import (
	"context"
	"errors"
	"testing"

	"github.com/lovely-eye/server/internal/auth"
	"github.com/stretchr/testify/require"
)

func TestRepositoryImplementsAuthBootstrapPolicyAtomically(t *testing.T) {
	t.Parallel()

	repo := New(setupTestDB(t))
	ctx := context.Background()
	first := &auth.StoredUser{Username: "admin", PasswordHash: "hash"}

	require.NoError(t, repo.CreateForRegistration(ctx, first, false))
	require.Equal(t, "admin", first.Role)
	require.NotZero(t, first.ID)
	require.False(t, first.CreatedAt.IsZero())

	err := repo.CreateForRegistration(ctx, &auth.StoredUser{
		Username:     "user",
		PasswordHash: "hash",
	}, false)
	require.ErrorIs(t, err, auth.ErrRegistrationDisabled)

	created, err := repo.CreateInitialAdminIfNoUsers(ctx, &auth.StoredUser{
		Username:     "ignored",
		PasswordHash: "hash",
	})
	require.NoError(t, err)
	require.False(t, created)
}

func TestRepositoryTranslatesMissingUsersToAuthError(t *testing.T) {
	t.Parallel()

	repo := New(setupTestDB(t))
	ctx := context.Background()

	_, err := repo.GetByID(ctx, 999)
	require.True(t, errors.Is(err, auth.ErrUserNotFound))
	_, err = repo.GetByUsername(ctx, "missing")
	require.True(t, errors.Is(err, auth.ErrUserNotFound))
}
