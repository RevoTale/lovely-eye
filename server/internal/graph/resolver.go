package graph

import (
	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthService      auth.Service
	SiteService      *services.SiteService
	AnalyticsService *services.AnalyticsService
	EventDefService  *services.EventDefinitionService
}

func NewResolver(
	authService auth.Service,
	siteService *services.SiteService,
	analyticsService *services.AnalyticsService,
	eventDefService *services.EventDefinitionService,
) *Resolver {
	return &Resolver{
		AuthService:      authService,
		SiteService:      siteService,
		AnalyticsService: analyticsService,
		EventDefService:  eventDefService,
	}
}
