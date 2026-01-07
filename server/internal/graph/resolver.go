package graph

import (
	"github.com/lovely-eye/server/internal/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthService      *services.AuthService
	SiteService      *services.SiteService
	AnalyticsService *services.AnalyticsService
}

func NewResolver(
	authService *services.AuthService,
	siteService *services.SiteService,
	analyticsService *services.AnalyticsService,
) *Resolver {
	return &Resolver{
		AuthService:      authService,
		SiteService:      siteService,
		AnalyticsService: analyticsService,
	}
}
