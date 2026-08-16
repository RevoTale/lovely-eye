package site_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	authpersistence "github.com/lovely-eye/server/internal/auth/persistence"
	"github.com/lovely-eye/server/internal/platform/database"
	"github.com/lovely-eye/server/internal/site"
	sitepersistence "github.com/lovely-eye/server/internal/site/persistence"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	_ "modernc.org/sqlite"
)

func TestSiteServicePropagatesDomainLookupFailure(t *testing.T) {
	service, db, userID := newSiteServiceTest(t)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	_, err := service.Create(context.Background(), site.CreateSiteInput{
		Domains: []string{"example.com"},
		Name:    "Example",
		UserID:  userID,
	})
	if err == nil {
		t.Fatal("expected storage failure")
	}
	if errors.Is(err, site.ErrSiteExists) {
		t.Fatalf("storage failure was classified as a domain conflict: %v", err)
	}
	if !strings.Contains(err.Error(), "check domain availability") {
		t.Fatalf("domain lookup failure was not preserved at its boundary: %v", err)
	}
}

func TestSiteServiceDistinguishesMissingSiteFromStorageFailure(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		service, _, userID := newSiteServiceTest(t)
		_, err := service.GetByID(context.Background(), 999, userID)
		if !errors.Is(err, site.ErrSiteNotFound) {
			t.Fatalf("expected site-not-found error, got %v", err)
		}
	})

	t.Run("storage failure", func(t *testing.T) {
		service, db, userID := newSiteServiceTest(t)
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
		_, err := service.GetByID(context.Background(), 999, userID)
		if err == nil {
			t.Fatal("expected storage failure")
		}
		if errors.Is(err, site.ErrSiteNotFound) {
			t.Fatalf("storage failure was classified as not found: %v", err)
		}
	})
}

func TestSiteServiceDomainConflictRemainsStable(t *testing.T) {
	service, _, userID := newSiteServiceTest(t)
	ctx := context.Background()
	_, err := service.Create(ctx, site.CreateSiteInput{
		Domains: []string{"example.com", "alias.example.com"},
		Name:    "First",
		UserID:  userID,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.Create(ctx, site.CreateSiteInput{
		Domains: []string{"example.com"},
		Name:    "Second",
		UserID:  userID,
	})
	if !errors.Is(err, site.ErrSiteExists) {
		t.Fatalf("expected stable domain-conflict error, got %v", err)
	}
}

func TestSiteServiceMissingMutationErrorsRemainStable(t *testing.T) {
	tests := map[string]func(context.Context, *site.Service, int64) error{
		"delete": func(ctx context.Context, service *site.Service, userID int64) error {
			return service.Delete(ctx, 999, userID)
		},
		"regenerate key": func(ctx context.Context, service *site.Service, userID int64) error {
			_, err := service.RegeneratePublicKey(ctx, 999, userID)
			return fmt.Errorf("regenerate key: %w", err)
		},
		"update": func(ctx context.Context, service *site.Service, userID int64) error {
			_, err := service.Update(ctx, 999, userID, site.UpdateSiteInput{Name: "Missing"})
			return fmt.Errorf("update: %w", err)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			service, _, userID := newSiteServiceTest(t)
			if err := run(context.Background(), service, userID); !errors.Is(err, site.ErrSiteNotFound) {
				t.Fatalf("expected stable site-not-found error, got %v", err)
			}
		})
	}
}

func newSiteServiceTest(t *testing.T) (*site.Service, *bun.DB, int64) {
	t.Helper()
	sqlDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	db := bun.NewDB(sqlDB, sqlitedialect.New())
	t.Cleanup(func() { _ = db.Close() })
	if err := database.Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	user := &authpersistence.User{Username: "site-test", PasswordHash: "hash", Role: "admin"}
	if _, err := db.NewInsert().Model(user).Exec(context.Background()); err != nil {
		t.Fatal(err)
	}
	return site.NewService(sitepersistence.New(db)), db, user.ID
}
