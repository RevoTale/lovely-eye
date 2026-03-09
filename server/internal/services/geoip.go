package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	geoipcore "github.com/lovely-eye/server/internal/geoip"
	"github.com/lovely-eye/server/internal/geoip/downloader"
	"github.com/lovely-eye/server/internal/geoip/lookup"
)

const (
	geoIPStateDisabled    = geoipcore.StateDisabled
	geoIPStateMissing     = geoipcore.StateMissing
	geoIPStateDownloading = geoipcore.StateDownloading
	geoIPStateReady       = geoipcore.StateReady
	geoIPStateError       = geoipcore.StateError
)

type GeoIPConfig struct {
	DBPath            string
	DownloadURL       string
	MaxMindLicenseKey string
}

type GeoIPSource string

const (
	GeoIPSourceUnknown     GeoIPSource = ""
	GeoIPSourceFile        GeoIPSource = "file"
	GeoIPSourceDownloadURL GeoIPSource = "download-url"
	GeoIPSourceDBIP        GeoIPSource = "dbip"
	GeoIPSourceMaxMind     GeoIPSource = "maxmind"
)

func (s GeoIPSource) String() string {
	return string(s)
}

type GeoIPStatus struct {
	State     string
	DBPath    string
	Source    GeoIPSource
	LastError string
	UpdatedAt *time.Time
}

type GeoIPCountry struct {
	Code string
	Name string
}

type Country struct {
	Name    string
	ISOCode string
}

var ErrNoDBReader = geoipcore.ErrNoDBReader

var UnknownCountry = Country{
	Name:    "Unknown",
	ISOCode: "-",
}

var LocalNetworkCountry = Country{
	Name:    "Local Network",
	ISOCode: "-",
}

type geoIPLookup interface {
	HasReader() bool
	FileExists() bool
	UpdatedAt() *time.Time
	Load() error
	ListCountries(search string) ([]geoipcore.ListedCountry, error)
	ResolveCountry(ipStr string) (geoipcore.Country, error)
	Close() error
}

type geoIPDownloader interface {
	HasDownloadSource() bool
	ConfiguredSource() geoipcore.Source
	BuildDownloadPlan() (downloader.DownloadPlan, error)
	Download(ctx context.Context, plan downloader.DownloadPlan) error
}

type GeoIPService struct {
	dbPath string

	lookup     geoIPLookup
	downloader geoIPDownloader

	status   GeoIPStatus
	statusMu sync.RWMutex

	enabled    bool
	downloadMu sync.Mutex
}

func NewGeoIPService(cfg GeoIPConfig) *GeoIPService {
	coreCfg := geoipcore.Config{
		DBPath:            cfg.DBPath,
		DownloadURL:       cfg.DownloadURL,
		MaxMindLicenseKey: cfg.MaxMindLicenseKey,
	}

	service := &GeoIPService{
		dbPath:     cfg.DBPath,
		lookup:     lookup.New(cfg.DBPath),
		downloader: downloader.New(coreCfg),
	}
	service.setStatus(GeoIPStatus{
		State:  geoIPStateDisabled,
		DBPath: cfg.DBPath,
	})
	return service
}

func (g *GeoIPService) SetEnabled(enabled bool) {
	g.statusMu.Lock()
	g.enabled = enabled
	g.statusMu.Unlock()

	if !enabled {
		_ = g.lookup.Close()
		g.setStatus(GeoIPStatus{
			State:  geoIPStateDisabled,
			DBPath: g.dbPath,
		})
	}
}

func (g *GeoIPService) Status() GeoIPStatus {
	g.statusMu.RLock()
	defer g.statusMu.RUnlock()
	return g.status
}

func (g *GeoIPService) ListCountries(search string) ([]GeoIPCountry, error) {
	status := g.Status()
	if status.State != geoIPStateReady {
		return []GeoIPCountry{}, nil
	}

	countries, err := g.lookup.ListCountries(search)
	if err != nil {
		return nil, fmt.Errorf("list GeoIP countries: %w", err)
	}
	return newGeoIPCountries(countries), nil
}

func (g *GeoIPService) EnsureAvailable(ctx context.Context) error {
	g.statusMu.RLock()
	enabled := g.enabled
	g.statusMu.RUnlock()

	if !enabled {
		g.setStatus(GeoIPStatus{
			State:  geoIPStateDisabled,
			DBPath: g.dbPath,
		})
		return nil
	}

	return g.loadDatabase(ctx, false)
}

func (g *GeoIPService) Refresh(ctx context.Context) error {
	return g.loadDatabase(ctx, true)
}

func (g *GeoIPService) ResolveCountry(ipStr string) (Country, error) {
	country, err := g.lookup.ResolveCountry(ipStr)
	return newCountry(country), err
}

func (g *GeoIPService) Close() error {
	if err := g.lookup.Close(); err != nil {
		return fmt.Errorf("close GeoIP lookup: %w", err)
	}
	return nil
}

func (g *GeoIPService) loadDatabase(ctx context.Context, forceRefresh bool) error {
	if g.dbPath == "" {
		return g.setStatusError(geoIPStateMissing, GeoIPSourceUnknown, errors.New("GeoIP database path is not configured"))
	}

	var loadErr error
	if !forceRefresh {
		if g.lookup.HasReader() {
			g.setReadyStatus(GeoIPSourceFile)
			return nil
		}

		if g.lookup.FileExists() {
			if err := g.lookup.Load(); err == nil {
				g.setReadyStatus(GeoIPSourceFile)
				return nil
			} else {
				loadErr = err
			}
		}
	}

	if !g.downloader.HasDownloadSource() {
		if loadErr != nil {
			return g.setStatusError(geoIPStateError, GeoIPSourceFile, loadErr)
		}
		return g.setStatusError(geoIPStateMissing, GeoIPSourceUnknown, errors.New("GeoIP download source is not configured"))
	}

	g.downloadMu.Lock()
	defer g.downloadMu.Unlock()

	if !forceRefresh && g.lookup.HasReader() {
		g.setReadyStatus(GeoIPSourceFile)
		return nil
	}

	plan, err := g.downloader.BuildDownloadPlan()
	if err != nil {
		return g.setStatusError(geoIPStateMissing, newGeoIPSource(g.downloader.ConfiguredSource()), err)
	}

	source := newGeoIPSource(plan.Source)
	g.setStatus(GeoIPStatus{
		State:  geoIPStateDownloading,
		DBPath: g.dbPath,
		Source: source,
	})

	if err := g.downloader.Download(ctx, plan); err != nil {
		return g.setStatusError(geoIPStateError, source, err)
	}

	if err := g.lookup.Load(); err != nil {
		return g.setStatusError(geoIPStateError, source, err)
	}

	g.setReadyStatus(source)
	return nil
}

func (g *GeoIPService) setReadyStatus(source GeoIPSource) {
	g.setStatus(GeoIPStatus{
		State:     geoIPStateReady,
		DBPath:    g.dbPath,
		Source:    source,
		UpdatedAt: g.lookup.UpdatedAt(),
	})
}

func (g *GeoIPService) setStatus(status GeoIPStatus) {
	if status.UpdatedAt == nil {
		now := time.Now()
		status.UpdatedAt = &now
	}
	g.statusMu.Lock()
	g.status = status
	g.statusMu.Unlock()
}

func (g *GeoIPService) setStatusError(state string, source GeoIPSource, err error) error {
	g.setStatus(GeoIPStatus{
		State:     state,
		DBPath:    g.dbPath,
		Source:    source,
		LastError: err.Error(),
	})
	return err
}

func newGeoIPSource(source geoipcore.Source) GeoIPSource {
	return GeoIPSource(source)
}

func newCountry(country geoipcore.Country) Country {
	return Country{
		Name:    country.Name,
		ISOCode: country.ISOCode,
	}
}

func newGeoIPCountries(countries []geoipcore.ListedCountry) []GeoIPCountry {
	result := make([]GeoIPCountry, 0, len(countries))
	for _, country := range countries {
		result = append(result, GeoIPCountry{
			Code: country.Code,
			Name: country.Name,
		})
	}
	return result
}
