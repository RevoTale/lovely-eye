package services

import (
	"context"
	"testing"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestCountryService_SyncFromGeoIP_PersistsNormalizedCountries(t *testing.T) {
	t.Parallel()

	db := setupAnalyticsServiceTestDB(t)
	countryRepo := repository.NewCountryRepository(db)
	geoIP := &fakeGeoIPProvider{
		countries: []GeoIPCountry{
			{Code: "us", Name: "United States"},
			{Code: "DE", Name: "Germany"},
			{Code: "-", Name: "Unknown"},
			{Code: "", Name: "Ignored"},
			{Code: "FR", Name: ""},
		},
	}
	service := NewCountryService(countryRepo, geoIP)

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

	db := setupAnalyticsServiceTestDB(t)
	countryRepo := repository.NewCountryRepository(db)
	service := NewCountryService(countryRepo, nil)
	ctx := context.Background()

	err := countryRepo.UpsertCountries(ctx, []models.Country{{Code: "US", Name: "United States"}})
	require.NoError(t, err)

	require.Equal(t, "United States", service.Name(ctx, "us"))
	require.Equal(t, "Unknown", service.Name(ctx, ""))
	require.Equal(t, "Unknown", service.Name(ctx, "-"))
	require.Equal(t, "ZZ", service.Name(ctx, "zz"))
}

func TestCountryService_List_ByCodeUsesRequestedOrder(t *testing.T) {
	t.Parallel()

	db := setupAnalyticsServiceTestDB(t)
	countryRepo := repository.NewCountryRepository(db)
	service := NewCountryService(countryRepo, nil)
	ctx := context.Background()

	err := countryRepo.UpsertCountries(ctx, []models.Country{
		{Code: "DE", Name: "Germany"},
		{Code: "US", Name: "United States"},
	})
	require.NoError(t, err)

	countries, err := service.List(ctx, "", []string{"us", "ZZ", "DE", "US"}, 10, 0)
	require.NoError(t, err)
	require.Equal(t, []CountryInfo{
		{Code: "US", Name: "United States"},
		{Code: "ZZ", Name: ""},
		{Code: "DE", Name: "Germany"},
	}, countries)
}
