package analytics

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	authpersistence "github.com/lovely-eye/server/internal/auth/persistence"
	"github.com/lovely-eye/server/internal/platform/database"
	sitepersistence "github.com/lovely-eye/server/internal/site/persistence"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	_ "modernc.org/sqlite"
)

type fakeGeoIPProvider struct {
	ensureCalls  int
	refreshCalls int
	enabled      bool

	ensureErr       error
	refreshErr      error
	resolveErr      error
	resolvedCountry Country

	countries []GeoIPCountry
	status    GeoIPStatus
}

func (f *fakeGeoIPProvider) SetEnabled(enabled bool) error {
	f.enabled = enabled
	return nil
}

func (f *fakeGeoIPProvider) Status() GeoIPStatus {
	if f.status.State != "" {
		return f.status
	}
	return GeoIPStatus{State: geoIPStateReady}
}

func (f *fakeGeoIPProvider) RecordFailure(err error) {
	f.status = GeoIPStatus{State: geoIPStateError, LastError: err.Error()}
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

func setupServiceTestDB(t *testing.T) *bun.DB {
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

func createGeoIPRequiredSite(t *testing.T, db *bun.DB) *sitepersistence.Site {
	t.Helper()

	ctx := context.Background()

	user := &authpersistence.User{
		Username:     "testuser",
		PasswordHash: "hash",
		Role:         "admin",
	}
	_, err := db.NewInsert().Model(user).Exec(ctx)
	require.NoError(t, err)

	site := &sitepersistence.Site{
		UserID:       user.ID,
		Name:         "Test Site",
		PublicKey:    "test-key",
		TrackCountry: true,
	}
	_, err = db.NewInsert().Model(site).Exec(ctx)
	require.NoError(t, err)

	return site
}

func TestService_SyncGeoIPRequirement_SyncsCountriesAfterEnsure(t *testing.T) {
	t.Parallel()

	db := setupServiceTestDB(t)
	createGeoIPRequiredSite(t, db)

	geoIP := &fakeGeoIPProvider{}
	countrySyncer := &fakeCountrySyncer{}

	service := NewService(
		analyticspersistence.New(db),
		sitepersistence.New(db),
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

func TestService_SyncGeoIPRequirement_NoSitesIsDisabledWithoutError(t *testing.T) {
	t.Parallel()

	db := setupServiceTestDB(t)
	geoIP := &fakeGeoIPProvider{}
	service := NewService(
		analyticspersistence.New(db),
		sitepersistence.New(db),
		nil,
		geoIP,
		nil,
		"test-analytics-identity-secret-32chars",
	)

	err := service.SyncGeoIPRequirement(context.Background())
	require.NoError(t, err)
	require.Zero(t, geoIP.ensureCalls)
	require.False(t, geoIP.enabled)
}

func TestService_RefreshGeoIPDatabase_SyncsCountriesAfterRefresh(t *testing.T) {
	t.Parallel()

	geoIP := &fakeGeoIPProvider{}
	countrySyncer := &fakeCountrySyncer{}

	service := NewService(nil, nil, nil, geoIP, countrySyncer, "test-analytics-identity-secret-32chars")

	status, err := service.RefreshGeoIPDatabase(context.Background())
	require.NoError(t, err)
	require.Equal(t, geoIPStateReady, status.State)
	require.Equal(t, 1, geoIP.refreshCalls)
	require.Equal(t, 1, countrySyncer.syncCalls)
}

func TestService_RefreshGeoIPDatabase_PropagatesCountrySyncError(t *testing.T) {
	t.Parallel()

	geoIP := &fakeGeoIPProvider{}
	countrySyncer := &fakeCountrySyncer{syncErr: errors.New("write failed")}

	service := NewService(nil, nil, nil, geoIP, countrySyncer, "test-analytics-identity-secret-32chars")

	_, err := service.RefreshGeoIPDatabase(context.Background())
	require.Error(t, err)
	require.ErrorContains(t, err, "sync persisted countries")
	require.ErrorContains(t, err, "write failed")
	require.Equal(t, geoIPStateError, geoIP.Status().State)
	require.Contains(t, geoIP.Status().LastError, "write failed")
}
