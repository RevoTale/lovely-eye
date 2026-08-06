package country_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/lovely-eye/server/internal/country"
	countrypersistence "github.com/lovely-eye/server/internal/country/persistence"
	"github.com/lovely-eye/server/internal/geoip"
	"github.com/lovely-eye/server/internal/platform/database"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	_ "modernc.org/sqlite"
)

func TestCountryService_SyncFromGeoIP_PersistsNormalizedCountries(t *testing.T) {
	t.Parallel()

	db := setupCountryDB(t)
	countryRepo := countrypersistence.New(db)
	geoIPProvider := &fakeGeoIPProvider{
		countries: []geoip.ListedCountry{
			{Code: "us", Name: "United States"},
			{Code: "DE", Name: "Germany"},
			{Code: "-", Name: "Unknown"},
			{Code: "", Name: "Ignored"},
			{Code: "FR", Name: ""},
		},
	}
	service := country.NewService(countryRepo, geoIPProvider)

	err := service.SyncFromGeoIP(context.Background())
	require.NoError(t, err)

	countries, err := countryRepo.GetCountriesByCodes(context.Background(), []string{"US", "DE", "FR"})
	require.NoError(t, err)
	require.Len(t, countries, 2)
	countryByCode := make(map[string]string, len(countries))
	for _, country := range countries {
		countryByCode[country.Code] = country.Name
	}
	require.Equal(t, map[string]string{
		"DE": "Germany",
		"US": "United States",
	}, countryByCode)
}

func TestCountryService_NameFallbacks(t *testing.T) {
	t.Parallel()

	db := setupCountryDB(t)
	countryRepo := countrypersistence.New(db)
	service := country.NewService(countryRepo, nil)
	ctx := context.Background()

	err := countryRepo.UpsertCountries(ctx, []country.Info{{Code: "US", Name: "United States"}})
	require.NoError(t, err)

	require.Equal(t, "United States", service.Name(ctx, "us"))
	require.Equal(t, "Unknown", service.Name(ctx, ""))
	require.Equal(t, "Unknown", service.Name(ctx, "-"))
	require.Equal(t, "ZZ", service.Name(ctx, "zz"))
}

func TestCountryService_List_ByCodeUsesRequestedOrder(t *testing.T) {
	t.Parallel()

	db := setupCountryDB(t)
	countryRepo := countrypersistence.New(db)
	service := country.NewService(countryRepo, nil)
	ctx := context.Background()

	err := countryRepo.UpsertCountries(ctx, []country.Info{
		{Code: "DE", Name: "Germany"},
		{Code: "US", Name: "United States"},
	})
	require.NoError(t, err)

	countries, err := service.List(ctx, "", []string{"us", "ZZ", "DE", "US"}, 10, 0)
	require.NoError(t, err)
	require.Equal(t, []country.Info{
		{Code: "US", Name: "United States"},
		{Code: "ZZ", Name: ""},
		{Code: "DE", Name: "Germany"},
	}, countries)
}

type fakeGeoIPProvider struct {
	countries []geoip.ListedCountry
}

func (f *fakeGeoIPProvider) ListCountries(string) ([]geoip.ListedCountry, error) {
	return f.countries, nil
}

func setupCountryDB(t *testing.T) *bun.DB {
	t.Helper()
	sqlDB, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	db := bun.NewDB(sqlDB, sqlitedialect.New())
	require.NoError(t, database.Migrate(context.Background(), db))
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return db
}
