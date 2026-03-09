package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/lovely-eye/server/internal/models"
)

type countryStore interface {
	UpsertCountries(ctx context.Context, countries []models.Country) error
	SearchCountries(ctx context.Context, search string, limit, offset int) ([]models.Country, error)
	GetCountriesByCodes(ctx context.Context, codes []string) ([]models.Country, error)
	GetCountryByCode(ctx context.Context, code string) (*models.Country, error)
}

type countryGeoIPProvider interface {
	ListCountries(search string) ([]GeoIPCountry, error)
}

type countrySyncer interface {
	SyncFromGeoIP(ctx context.Context) error
}

type CountryInfo struct {
	Code string
	Name string
}

type CountryService struct {
	countryRepo  countryStore
	geoIPService countryGeoIPProvider
}

func NewCountryService(countryRepo countryStore, geoIPService countryGeoIPProvider) *CountryService {
	return &CountryService{
		countryRepo:  countryRepo,
		geoIPService: geoIPService,
	}
}

func (s *CountryService) SyncFromGeoIP(ctx context.Context) error {
	if s.geoIPService == nil || s.countryRepo == nil {
		return nil
	}

	countries, err := s.geoIPService.ListCountries("")
	if err != nil {
		return fmt.Errorf("list geoip countries: %w", err)
	}

	persistedCountries := make([]models.Country, 0, len(countries))
	for _, country := range countries {
		code := normalizeCountryCode(country.Code)
		name := strings.TrimSpace(country.Name)
		if code == "" || code == "-" || name == "" {
			continue
		}

		persistedCountries = append(persistedCountries, models.Country{
			Code: code,
			Name: name,
		})
	}

	if err := s.countryRepo.UpsertCountries(ctx, persistedCountries); err != nil {
		return fmt.Errorf("upsert persisted countries: %w", err)
	}

	return nil
}

func (s *CountryService) List(ctx context.Context, search string, codes []string, limit, offset int) ([]CountryInfo, error) {
	if s.countryRepo == nil {
		return []CountryInfo{}, errors.New("country repository is nil")
	}

	if len(codes) > 0 {
		return s.lookupByCode(ctx, codes, limit, offset)
	}

	countries, err := s.countryRepo.SearchCountries(ctx, search, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search countries: %w", err)
	}

	result := make([]CountryInfo, 0, len(countries))
	for _, country := range countries {
		result = append(result, CountryInfo{
			Code: country.Code,
			Name: country.Name,
		})
	}
	return result, nil
}

func (s *CountryService) Name(ctx context.Context, code string) string {
	normalizedCode := normalizeCountryCode(code)
	switch normalizedCode {
	case "", "-":
		return "Unknown"
	}

	if s.countryRepo == nil {
		return normalizedCode
	}

	country, err := s.countryRepo.GetCountryByCode(ctx, normalizedCode)
	if err != nil {
		slog.Error("country lookup failed", "code", normalizedCode, "error", err)
		return normalizedCode
	}
	if country == nil || strings.TrimSpace(country.Name) == "" {
		return normalizedCode
	}

	return country.Name
}

func (s *CountryService) lookupByCode(ctx context.Context, codes []string, limit, offset int) ([]CountryInfo, error) {
	if len(codes) == 0 {
		return []CountryInfo{}, nil
	}

	normalizedCodes := make([]string, 0, len(codes))
	seen := make(map[string]struct{}, len(codes))
	for _, code := range codes {
		normalizedCode := normalizeCountryCode(code)
		if normalizedCode == "" {
			continue
		}
		if _, ok := seen[normalizedCode]; ok {
			continue
		}
		seen[normalizedCode] = struct{}{}
		normalizedCodes = append(normalizedCodes, normalizedCode)
	}

	if offset >= len(normalizedCodes) {
		return []CountryInfo{}, nil
	}

	end := len(normalizedCodes)
	if limit > 0 && offset+limit < end {
		end = offset + limit
	}

	requestedCodes := normalizedCodes[offset:end]
	countries, err := s.countryRepo.GetCountriesByCodes(ctx, requestedCodes)
	if err != nil {
		return nil, fmt.Errorf("get countries by code: %w", err)
	}

	countryByCode := make(map[string]string, len(countries))
	for _, country := range countries {
		countryByCode[country.Code] = country.Name
	}

	result := make([]CountryInfo, 0, len(requestedCodes))
	for _, code := range requestedCodes {
		result = append(result, CountryInfo{
			Code: code,
			Name: countryByCode[code],
		})
	}

	return result, nil
}

func normalizeCountryCode(code string) string {
	trimmedCode := strings.TrimSpace(code)
	if trimmedCode == "" || trimmedCode == "-" {
		return trimmedCode
	}

	return strings.ToUpper(trimmedCode)
}
