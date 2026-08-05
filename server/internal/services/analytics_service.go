package services

import (
	"context"
	"time"

	"github.com/lovely-eye/server/internal/repository"
)

const (
	activeSessionWindow              = 30 * time.Minute
	defaultMaxSinglePageExitDuration = 4 * time.Hour
)

type AnalyticsService struct {
	analyticsRepo         *repository.AnalyticsRepository
	countryService        countrySyncer
	siteRepo              *repository.SiteRepository
	eventDefinitionRepo   *repository.EventDefinitionRepository
	botDetector           *BotDetector
	geoIPService          geoIPProvider
	identitySecret        []byte
	maxSinglePageDuration time.Duration
	now                   func() time.Time
}

type geoIPProvider interface {
	SetEnabled(enabled bool)
	Status() GeoIPStatus
	EnsureAvailable(ctx context.Context) error
	Refresh(ctx context.Context) error
	ResolveCountry(ipStr string) (Country, error)
	ListCountries(search string) ([]GeoIPCountry, error)
	Close() error
}

func NewAnalyticsService(
	analyticsRepo *repository.AnalyticsRepository,
	siteRepo *repository.SiteRepository,
	eventDefinitionRepo *repository.EventDefinitionRepository,
	geoIPService geoIPProvider,
	countryService countrySyncer,
	identitySecret string,
) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo:         analyticsRepo,
		countryService:        countryService,
		siteRepo:              siteRepo,
		eventDefinitionRepo:   eventDefinitionRepo,
		botDetector:           NewBotDetector(),
		geoIPService:          geoIPService,
		identitySecret:        []byte(identitySecret),
		maxSinglePageDuration: defaultMaxSinglePageExitDuration,
		now:                   time.Now,
	}
}

func (s *AnalyticsService) SetMaxSinglePageDuration(duration time.Duration) {
	if duration <= 0 {
		duration = defaultMaxSinglePageExitDuration
	}
	s.maxSinglePageDuration = duration
}

func (s *AnalyticsService) sessionLookupWindow(exit bool) time.Duration {
	if !exit {
		return activeSessionWindow
	}
	if s.maxSinglePageDuration <= 0 {
		return defaultMaxSinglePageExitDuration
	}
	return s.maxSinglePageDuration
}

type CollectInput struct {
	SiteKey     string
	Path        string
	Exit        bool
	Referrer    string
	UserAgent   string
	IP          string
	Origin      string
	Referer     string
	UTMSource   string
	UTMMedium   string
	UTMCampaign string
}

type EventInput struct {
	SiteKey    string
	Name       string
	Path       string
	Properties string
	UserAgent  string
	IP         string
	Origin     string
	Referer    string
}
