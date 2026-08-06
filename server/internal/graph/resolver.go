package graph

import (
	"context"
	"net/http"

	"github.com/lovely-eye/server/internal/analytics"
	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/country"
	"github.com/lovely-eye/server/internal/event"
	"github.com/lovely-eye/server/internal/site"
)

type AuthService interface {
	Register(ctx context.Context, input auth.RegisterInput) (*auth.User, *auth.Tokens, error)
	Login(ctx context.Context, input auth.LoginInput) (*auth.User, *auth.Tokens, error)
	GetUserByID(ctx context.Context, id int64) (*auth.User, error)
	RegistrationStatus(ctx context.Context) (*auth.RegistrationStatus, error)
}

type AuthCookies interface {
	SetAuthCookies(w http.ResponseWriter, tokens *auth.Tokens)
	ClearAuthCookies(w http.ResponseWriter)
}

type Resolver struct {
	AuthService      AuthService
	AuthCookies      AuthCookies
	SiteService      *site.Service
	AnalyticsService *analytics.Service
	CountryService   *country.Service
	EventDefService  *event.Service
	DashboardLimits  DashboardLimits
}

type DashboardLimits struct {
	MaxDailyRangeDays     int
	MaxHourlyRangeDays    int
	MaxFilterValues       int
	MaxFilterStringLength int
}

func NewResolver(
	authService AuthService,
	authCookies AuthCookies,
	siteService *site.Service,
	analyticsService *analytics.Service,
	countryService *country.Service,
	eventDefService *event.Service,
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
		AuthCookies:      authCookies,
		SiteService:      siteService,
		AnalyticsService: analyticsService,
		CountryService:   countryService,
		EventDefService:  eventDefService,
		DashboardLimits:  dashboardLimits,
	}
}
