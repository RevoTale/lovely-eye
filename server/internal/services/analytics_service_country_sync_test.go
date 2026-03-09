package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/lovely-eye/server/internal/database"
	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	_ "modernc.org/sqlite"
)

type fakeGeoIPProvider struct {
	ensureCalls  int
	refreshCalls int

	ensureErr       error
	refreshErr      error
	resolveErr      error
	resolvedCountry Country

	countries []GeoIPCountry
}

func (f *fakeGeoIPProvider) SetEnabled(bool) {}

func (f *fakeGeoIPProvider) Status() GeoIPStatus {
	return GeoIPStatus{State: geoIPStateReady}
}

func (f *fakeGeoIPProvider) EnsureAvailable(context.Context) error {
	f.ensureCalls++
	return f.ensureErr
}

func (f *fakeGeoIPProvider) Refresh(context.Context) error {
	f.refreshCalls++
	return f.refreshErr
}

func (f *fakeGeoIPProvider) ResolveCountry(string) (Country, error) {
	if f.resolveErr != nil {
		return Country{}, f.resolveErr
	}
	if f.resolvedCountry != (Country{}) {
		return f.resolvedCountry, nil
	}
	return UnknownCountry, nil
}

func (f *fakeGeoIPProvider) ListCountries(string) ([]GeoIPCountry, error) {
	return f.countries, nil
}

func (f *fakeGeoIPProvider) Close() error {
	return nil
}

type fakeCountrySyncer struct {
	syncCalls int
	syncErr   error
}

func (f *fakeCountrySyncer) SyncFromGeoIP(context.Context) error {
	f.syncCalls++
	return f.syncErr
}

func setupAnalyticsServiceTestDB(t *testing.T) *bun.DB {
	t.Helper()

	sqldb, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	db := bun.NewDB(sqldb, sqlitedialect.New())
	err = database.Migrate(context.Background(), db)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	return db
}

func createGeoIPRequiredSite(t *testing.T, db *bun.DB) *models.Site {
	t.Helper()

	ctx := context.Background()

	user := &models.User{
		Username:     "testuser",
		PasswordHash: "hash",
		Role:         "admin",
	}
	_, err := db.NewInsert().Model(user).Exec(ctx)
	require.NoError(t, err)

	site := &models.Site{
		UserID:       user.ID,
		Name:         "Test Site",
		PublicKey:    "test-key",
		TrackCountry: true,
	}
	_, err = db.NewInsert().Model(site).Exec(ctx)
	require.NoError(t, err)

	return site
}

func TestAnalyticsService_SyncGeoIPRequirement_SyncsCountriesAfterEnsure(t *testing.T) {
	t.Parallel()

	db := setupAnalyticsServiceTestDB(t)
	createGeoIPRequiredSite(t, db)

	geoIP := &fakeGeoIPProvider{}
	countrySyncer := &fakeCountrySyncer{}

	service := NewAnalyticsService(
		repository.NewAnalyticsRepository(db),
		repository.NewSiteRepository(db),
		nil,
		geoIP,
		countrySyncer,
		"test-analytics-identity-secret-32chars",
	)

	err := service.SyncGeoIPRequirement(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, geoIP.ensureCalls)
	require.Equal(t, 1, countrySyncer.syncCalls)
}

func TestAnalyticsService_RefreshGeoIPDatabase_SyncsCountriesAfterRefresh(t *testing.T) {
	t.Parallel()

	geoIP := &fakeGeoIPProvider{}
	countrySyncer := &fakeCountrySyncer{}

	service := NewAnalyticsService(nil, nil, nil, geoIP, countrySyncer, "test-analytics-identity-secret-32chars")

	status, err := service.RefreshGeoIPDatabase(context.Background())
	require.NoError(t, err)
	require.Equal(t, geoIPStateReady, status.State)
	require.Equal(t, 1, geoIP.refreshCalls)
	require.Equal(t, 1, countrySyncer.syncCalls)
}

func TestAnalyticsService_RefreshGeoIPDatabase_PropagatesCountrySyncError(t *testing.T) {
	t.Parallel()

	geoIP := &fakeGeoIPProvider{}
	countrySyncer := &fakeCountrySyncer{syncErr: errors.New("write failed")}

	service := NewAnalyticsService(nil, nil, nil, geoIP, countrySyncer, "test-analytics-identity-secret-32chars")

	_, err := service.RefreshGeoIPDatabase(context.Background())
	require.Error(t, err)
	require.ErrorContains(t, err, "sync persisted countries")
	require.ErrorContains(t, err, "write failed")
}
