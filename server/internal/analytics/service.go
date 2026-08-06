package analytics

import (
	"context"
	"time"

	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	"github.com/lovely-eye/server/internal/event"
	"github.com/lovely-eye/server/internal/site"
)

const (
	activeSessionWindow              = 30 * time.Minute
	defaultMaxSinglePageExitDuration = 4 * time.Hour
)

type Service struct {
	analyticsRepo         *analyticspersistence.Repository
	countryService        countrySyncer
	siteRepo              site.Store
	eventDefinitionStore  event.Store
	botDetector           *BotDetector
	geoIPService          geoIPProvider
	identitySecret        []byte
	maxSinglePageDuration time.Duration
	now                   func() time.Time
}

type countrySyncer interface {
	SyncFromGeoIP(ctx context.Context) error
}

type geoIPProvider interface {
	SetEnabled(enabled bool) error
	Status() GeoIPStatus
	RecordFailure(err error)
	EnsureAvailable(ctx context.Context) error
	Refresh(ctx context.Context) error
	ResolveCountry(ipStr string) (Country, error)
	ListCountries(search string) ([]GeoIPCountry, error)
	Close() error
}

func NewService(
	analyticsRepo *analyticspersistence.Repository,
	siteRepo site.Store,
	eventDefinitionStore event.Store,
	geoIPService geoIPProvider,
	countryService countrySyncer,
	identitySecret string,
) *Service {
	return &Service{
		analyticsRepo:         analyticsRepo,
		countryService:        countryService,
		siteRepo:              siteRepo,
		eventDefinitionStore:  eventDefinitionStore,
		botDetector:           NewBotDetector(),
		geoIPService:          geoIPService,
		identitySecret:        []byte(identitySecret),
		maxSinglePageDuration: defaultMaxSinglePageExitDuration,
		now:                   time.Now,
	}
}

func (s *Service) SetMaxSinglePageDuration(duration time.Duration) {
	if duration <= 0 {
		duration = defaultMaxSinglePageExitDuration
	}
	s.maxSinglePageDuration = duration
}

func (s *Service) sessionLookupWindow(exit bool) time.Duration {
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
