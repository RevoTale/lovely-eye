package graph

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/site"
	"github.com/stretchr/testify/require"
)

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want errorCode
	}{
		{name: "resolver authentication", err: unauthenticated(), want: errorCodeUnauthenticated},
		{name: "feature authorization", err: fmt.Errorf("load site: %w", site.ErrNotAuthorized), want: errorCodeForbidden},
		{name: "missing resource", err: fmt.Errorf("load site: %w", site.ErrSiteNotFound), want: errorCodeNotFound},
		{name: "uniqueness conflict", err: fmt.Errorf("create user: %w", auth.ErrUserExists), want: errorCodeConflict},
		{name: "feature validation", err: fmt.Errorf("validate site: %w", site.ErrInvalidDomain), want: errorCodeBadUserInput},
		{name: "unknown failure", err: errors.New("database unavailable"), want: errorCodeInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, classifyError(tt.err))
		})
	}
}

func TestPresentErrorMasksInternalDetails(t *testing.T) {
	presented := presentError(context.Background(), errors.New("database password leaked"))

	require.Equal(t, "internal server error", presented.Message)
	require.Equal(t, string(errorCodeInternal), presented.Extensions["code"])
}
