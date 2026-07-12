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
	DashboardLimits  DashboardLimits
}

type DashboardLimits struct {
	MaxDailyRangeDays     int
	MaxHourlyRangeDays    int
	MaxFilterValues       int
	MaxFilterStringLength int
}

func NewResolver(
	authService auth.Service,
	siteService *services.SiteService,
	analyticsService *services.AnalyticsService,
	countryService *services.CountryService,
	eventDefService *services.EventDefinitionService,
	dashboardLimits DashboardLimits,
) *Resolver {
	if dashboardLimits.MaxDailyRangeDays <= 0 {
		dashboardLimits.MaxDailyRangeDays = 730
	}
	if dashboardLimits.MaxHourlyRangeDays <= 0 {
		dashboardLimits.MaxHourlyRangeDays = 31
	}
	if dashboardLimits.MaxFilterValues <= 0 {
		dashboardLimits.MaxFilterValues = 100
	}
	if dashboardLimits.MaxFilterStringLength <= 0 {
		dashboardLimits.MaxFilterStringLength = 2048
	}
	return &Resolver{
		AuthService:      authService,
		SiteService:      siteService,
		AnalyticsService: analyticsService,
		CountryService:   countryService,
		EventDefService:  eventDefService,
		DashboardLimits:  dashboardLimits,
	}
}
