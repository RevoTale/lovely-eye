package repository

import (
	"context"
	"testing"

	"github.com/lovely-eye/server/internal/models"
	"github.com/stretchr/testify/require"
)

func TestCountryRepository_UpsertCountries_InsertsUpdatesAndKeepsMissingRows(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	repo := NewCountryRepository(db)
	ctx := context.Background()

	err := repo.UpsertCountries(ctx, []models.Country{
		{Code: "US", Name: "United States"},
		{Code: "DE", Name: "Germany"},
	})
	require.NoError(t, err)

	err = repo.UpsertCountries(ctx, []models.Country{
		{Code: "US", Name: "USA"},
		{Code: "FR", Name: "France"},
	})
	require.NoError(t, err)

	var countries []models.Country
	err = db.NewSelect().
		Model(&countries).
		Order("co.code ASC").
		Scan(ctx)
	require.NoError(t, err)

	require.Equal(t, []models.Country{
		{Code: "DE", Name: "Germany"},
		{Code: "FR", Name: "France"},
		{Code: "US", Name: "USA"},
	}, countries)
}

func TestCountryRepository_SearchCountries(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	repo := NewCountryRepository(db)
	ctx := context.Background()

	err := repo.UpsertCountries(ctx, []models.Country{
		{Code: "US", Name: "United States"},
		{Code: "GB", Name: "United Kingdom"},
		{Code: "DE", Name: "Germany"},
	})
	require.NoError(t, err)

	countries, err := repo.SearchCountries(ctx, "uni", 10, 0)
	require.NoError(t, err)
	require.Len(t, countries, 2)
	require.Equal(t, "GB", countries[0].Code)
	require.Equal(t, "US", countries[1].Code)

	countries, err = repo.SearchCountries(ctx, "de", 10, 0)
	require.NoError(t, err)
	require.Len(t, countries, 1)
	require.Equal(t, "DE", countries[0].Code)
}

func TestCountryRepository_GetCountriesByCodes(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	repo := NewCountryRepository(db)
	ctx := context.Background()

	err := repo.UpsertCountries(ctx, []models.Country{
		{Code: "US", Name: "United States"},
		{Code: "DE", Name: "Germany"},
	})
	require.NoError(t, err)

	countries, err := repo.GetCountriesByCodes(ctx, []string{"DE", "US"})
	require.NoError(t, err)
	require.Len(t, countries, 2)

	countryByCode := make(map[string]string, len(countries))
	for _, country := range countries {
		countryByCode[country.Code] = country.Name
	}

	require.Equal(t, "Germany", countryByCode["DE"])
	require.Equal(t, "United States", countryByCode["US"])

	country, err := repo.GetCountryByCode(ctx, "US")
	require.NoError(t, err)
	require.NotNil(t, country)
	require.Equal(t, "United States", country.Name)

	country, err = repo.GetCountryByCode(ctx, "FR")
	require.NoError(t, err)
	require.Nil(t, country)
}
