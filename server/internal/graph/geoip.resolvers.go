package graph

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/graph/model"
)

// Name is the resolver for the name field.
func (r *countryResolver) Name(ctx context.Context, obj *model.Country) (string, error) {
	if obj.NameCache != nil && strings.TrimSpace(*obj.NameCache) != "" {
		return *obj.NameCache, nil
	}
	if r.CountryService == nil {
		return obj.Code, nil
	}
	return r.CountryService.Name(ctx, obj.Code), nil
}

// RefreshGeoIPDatabase is the resolver for the refreshGeoIPDatabase field.
func (r *mutationResolver) RefreshGeoIPDatabase(ctx context.Context) (*model.GeoIPStatus, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	status, err := r.AnalyticsService.RefreshGeoIPDatabase(ctx)
	if err != nil {
		// Refresh is an operational status mutation: return its actionable status instead of null GraphQL data.
		slog.WarnContext(ctx, "GeoIP refresh failed", "error", err)
	}
	return convertToGraphQLGeoIPStatus(status), nil
}

// GeoIPStatus is the resolver for the geoIPStatus field.
func (r *queryResolver) GeoIPStatus(ctx context.Context) (*model.GeoIPStatus, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	status := r.AnalyticsService.GeoIPStatus()
	return convertToGraphQLGeoIPStatus(status), nil
}

// GeoIPCountries is the resolver for the geoIPCountries field.
func (r *queryResolver) GeoIPCountries(ctx context.Context, search *string, codes []string, paging model.PagingInput) ([]*model.Country, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	query := ""
	if search != nil {
		query = *search
	}

	limit, offset := normalizePaging(paging)

	if r.CountryService == nil {
		return []*model.Country{}, nil
	}

	countries, err := r.CountryService.List(ctx, query, codes, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get geoip countries: %w", err)
	}

	result := make([]*model.Country, 0, len(countries))
	for _, country := range countries {
		result = append(result, newGraphQLCountry(country.Code, country.Name))
	}
	return result, nil
}

// Country returns CountryResolver implementation.
func (r *Resolver) Country() CountryResolver { return &countryResolver{r} }

type countryResolver struct{ *Resolver }
