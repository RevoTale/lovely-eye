package country

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/lovely-eye/server/internal/geoip"
)

var ErrNotFound = errors.New("country not found")

type Store interface {
	UpsertCountries(ctx context.Context, countries []Info) error
	SearchCountries(ctx context.Context, search string, limit, offset int) ([]Info, error)
	GetCountriesByCodes(ctx context.Context, codes []string) ([]Info, error)
	GetCountryByCode(ctx context.Context, code string) (*Info, error)
}

type countryGeoIPProvider interface {
	ListCountries(search string) ([]geoip.ListedCountry, error)
}

type Info struct {
	Code string
	Name string
}

type Service struct {
	countryRepo  Store
	geoIPService countryGeoIPProvider
}

func NewService(countryRepo Store, geoIPService countryGeoIPProvider) *Service {
	return &Service{
		countryRepo:  countryRepo,
		geoIPService: geoIPService,
	}
}

func (s *Service) SyncFromGeoIP(ctx context.Context) error {
	if s.geoIPService == nil || s.countryRepo == nil {
		return nil
	}

	countries, err := s.geoIPService.ListCountries("")
	if err != nil {
		return fmt.Errorf("list geoip countries: %w", err)
	}

	persistedCountries := make([]Info, 0, len(countries))
	for _, country := range countries {
		code := normalizeCountryCode(country.Code)
		name := strings.TrimSpace(country.Name)
		if code == "" || code == "-" || name == "" {
			continue
		}

		persistedCountries = append(persistedCountries, Info{
			Code: code,
			Name: name,
		})
	}

	if err := s.countryRepo.UpsertCountries(ctx, persistedCountries); err != nil {
		return fmt.Errorf("upsert persisted countries: %w", err)
	}

	return nil
}

func (s *Service) List(ctx context.Context, search string, codes []string, limit, offset int) ([]Info, error) {
	if s.countryRepo == nil {
		return []Info{}, errors.New("country repository is nil")
	}

	if len(codes) > 0 {
		return s.lookupByCode(ctx, codes, limit, offset)
	}

	countries, err := s.countryRepo.SearchCountries(ctx, search, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search countries: %w", err)
	}

	result := make([]Info, 0, len(countries))
	for _, country := range countries {
		result = append(result, Info{
			Code: country.Code,
			Name: country.Name,
		})
	}
	return result, nil
}

func (s *Service) Name(ctx context.Context, code string) string {
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
		if errors.Is(err, ErrNotFound) {
			return normalizedCode
		}
		slog.Error("country lookup failed", "code", normalizedCode, "error", err)
		return normalizedCode
	}
	if country == nil || strings.TrimSpace(country.Name) == "" {
		return normalizedCode
	}

	return country.Name
}

func (s *Service) lookupByCode(ctx context.Context, codes []string, limit, offset int) ([]Info, error) {
	if len(codes) == 0 {
		return []Info{}, nil
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
		return []Info{}, nil
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

	result := make([]Info, 0, len(requestedCodes))
	for _, code := range requestedCodes {
		result = append(result, Info{
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
