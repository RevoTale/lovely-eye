package graph

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/graph/model"
	"github.com/lovely-eye/server/internal/site"
)

// CreateSite is the resolver for the createSite field.
func (r *mutationResolver) CreateSite(ctx context.Context, input model.CreateSiteInput) (*model.Site, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	site, err := r.SiteService.Create(ctx, site.CreateSiteInput{
		Domains: input.Domains,
		Name:    input.Name,
		UserID:  claims.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create site: %w", err)
	}

	return buildGraphQLSite(site), nil
}

// UpdateSite is the resolver for the updateSite field.
func (r *mutationResolver) UpdateSite(ctx context.Context, id string, input model.UpdateSiteInput) (*model.Site, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	siteID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, badUserInput("invalid site ID")
	}

	site, err := r.SiteService.Update(ctx, siteID, claims.UserID, site.UpdateSiteInput{
		Name:             input.Name,
		TrackCountry:     input.TrackCountry,
		Domains:          input.Domains,
		BlockedIPs:       input.BlockedIPs,
		BlockedCountries: input.BlockedCountries,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update site: %w", err)
	}
	if err := r.AnalyticsService.SyncGeoIPRequirement(ctx); err != nil {
		// The site update is already durable; GeoIP status carries this best-effort follow-up failure to the UI.
		slog.WarnContext(ctx, "GeoIP synchronization failed after site update", "error", err)
	}

	return buildGraphQLSite(site), nil
}

// DeleteSite is the resolver for the deleteSite field.
func (r *mutationResolver) DeleteSite(ctx context.Context, id string) (bool, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return false, unauthenticated()
	}

	siteID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return false, badUserInput("invalid site ID")
	}

	if err := r.SiteService.Delete(ctx, siteID, claims.UserID); err != nil {
		return false, fmt.Errorf("failed to delete site: %w", err)
	}

	return true, nil
}

// RegenerateSiteKey is the resolver for the regenerateSiteKey field.
func (r *mutationResolver) RegenerateSiteKey(ctx context.Context, id string) (*model.Site, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	siteID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, badUserInput("invalid site ID")
	}

	site, err := r.SiteService.RegeneratePublicKey(ctx, siteID, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to regenerate site key: %w", err)
	}

	return buildGraphQLSite(site), nil
}

// Sites is the resolver for the sites field.
func (r *queryResolver) Sites(ctx context.Context, paging model.PagingInput) ([]*model.Site, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	limit, offset := normalizePaging(paging)
	sites, err := r.SiteService.GetUserSites(ctx, claims.UserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get sites: %w", err)
	}

	var result []*model.Site
	for _, site := range sites {
		result = append(result, buildGraphQLSite(site))
	}

	return result, nil
}

// Site is the resolver for the site field.
func (r *queryResolver) Site(ctx context.Context, id string) (*model.Site, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	siteID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, badUserInput("invalid site ID")
	}

	site, err := r.SiteService.GetByID(ctx, siteID, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	return buildGraphQLSite(site), nil
}
