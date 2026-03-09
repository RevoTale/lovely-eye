package graph

import (
	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/services"
)

type Resolver struct {
	AuthService      auth.Service
	SiteService      *services.SiteService
	AnalyticsService *services.AnalyticsService
	CountryService   *services.CountryService
	EventDefService  *services.EventDefinitionService
}

func NewResolver(
	authService auth.Service,
	siteService *services.SiteService,
	analyticsService *services.AnalyticsService,
	countryService *services.CountryService,
	eventDefService *services.EventDefinitionService,
) *Resolver {
	return &Resolver{
		AuthService:      authService,
		SiteService:      siteService,
		AnalyticsService: analyticsService,
		CountryService:   countryService,
		EventDefService:  eventDefService,
	}
}
