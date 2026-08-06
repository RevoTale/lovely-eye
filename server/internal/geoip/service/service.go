package service

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

type Config = geoipcore.Config
type Status = geoipcore.Status
type Source = geoipcore.Source
type State = geoipcore.State
type ListedCountry = geoipcore.ListedCountry
type Country = geoipcore.Country

var ErrNoDBReader = geoipcore.ErrNoDBReader
var UnknownCountry = geoipcore.UnknownCountry

const (
	StateDisabled    = geoipcore.StateDisabled
	StateMissing     = geoipcore.StateMissing
	StateDownloading = geoipcore.StateDownloading
	StateReady       = geoipcore.StateReady
	StateError       = geoipcore.StateError
	SourceUnknown    = geoipcore.SourceUnknown
	SourceFile       = geoipcore.SourceFile
	SourceDBIP       = geoipcore.SourceDBIP
	SourceMaxMind    = geoipcore.SourceMaxMind
)

type geoIPLookup interface {
	HasReader() bool
	FileExists() bool
	UpdatedAt() *time.Time
	Load() error
	ListCountries(search string) ([]ListedCountry, error)
	ResolveCountry(ipStr string) (Country, error)
	Close() error
}

type geoIPDownloader interface {
	HasDownloadSource() bool
	ConfiguredSource() Source
	BuildDownloadPlan() (downloader.DownloadPlan, error)
	Download(ctx context.Context, plan downloader.DownloadPlan) error
}

type Service struct {
	dbPath string

	lookup     geoIPLookup
	downloader geoIPDownloader

	status   Status
	statusMu sync.RWMutex

	enabled    bool
	downloadMu sync.Mutex
}

func NewService(cfg Config) *Service {
	service := &Service{
		dbPath:     cfg.DBPath,
		lookup:     lookup.New(cfg.DBPath),
		downloader: downloader.New(cfg),
	}
	service.setStatus(Status{
		State:  StateDisabled,
		DBPath: cfg.DBPath,
	})
	return service
}

func (g *Service) SetEnabled(enabled bool) error {
	g.statusMu.Lock()
	g.enabled = enabled
	g.statusMu.Unlock()

	if !enabled {
		if err := g.lookup.Close(); err != nil {
			failure := fmt.Errorf("close GeoIP lookup while disabling: %w", err)
			g.RecordFailure(failure)
			return failure
		}
		g.setStatus(Status{
			State:  StateDisabled,
			DBPath: g.dbPath,
		})
	}
	return nil
}

func (g *Service) Status() Status {
	g.statusMu.RLock()
	defer g.statusMu.RUnlock()
	return g.status
}

func (g *Service) RecordFailure(err error) {
	status := g.Status()
	if status.State != StateMissing {
		status.State = StateError
	}
	status.LastError = err.Error()
	status.UpdatedAt = nil
	g.setStatus(status)
}

func (g *Service) ListCountries(search string) ([]ListedCountry, error) {
	status := g.Status()
	if status.State != StateReady {
		return []ListedCountry{}, nil
	}

	countries, err := g.lookup.ListCountries(search)
	if err != nil {
		return nil, fmt.Errorf("list GeoIP countries: %w", err)
	}
	return countries, nil
}

func (g *Service) EnsureAvailable(ctx context.Context) error {
	g.statusMu.RLock()
	enabled := g.enabled
	g.statusMu.RUnlock()

	if !enabled {
		g.setStatus(Status{
			State:  StateDisabled,
			DBPath: g.dbPath,
		})
		return nil
	}

	return g.loadDatabase(ctx, false)
}

func (g *Service) Refresh(ctx context.Context) error {
	return g.loadDatabase(ctx, true)
}

func (g *Service) ResolveCountry(ipStr string) (Country, error) {
	country, err := g.lookup.ResolveCountry(ipStr)
	if err != nil {
		return country, fmt.Errorf("resolve country: %w", err)
	}
	return country, nil
}

func (g *Service) Close() error {
	if err := g.lookup.Close(); err != nil {
		return fmt.Errorf("close GeoIP lookup: %w", err)
	}
	return nil
}

func (g *Service) loadDatabase(ctx context.Context, forceRefresh bool) error {
	if g.dbPath == "" {
		return g.setStatusError(StateMissing, SourceUnknown, errors.New("GeoIP database path is not configured"))
	}

	var loadErr error
	if !forceRefresh {
		if g.lookup.HasReader() {
			g.setReadyStatus(SourceFile)
			return nil
		}

		if g.lookup.FileExists() {
			err := g.lookup.Load()
			if err == nil {
				g.setReadyStatus(SourceFile)
				return nil
			}
			loadErr = err
		}
	}

	if !g.downloader.HasDownloadSource() {
		if loadErr != nil {
			return g.setStatusError(StateError, SourceFile, loadErr)
		}
		return g.setStatusError(StateMissing, SourceUnknown, errors.New("GeoIP download source is not configured"))
	}

	g.downloadMu.Lock()
	defer g.downloadMu.Unlock()

	if !forceRefresh && g.lookup.HasReader() {
		g.setReadyStatus(SourceFile)
		return nil
	}

	plan, err := g.downloader.BuildDownloadPlan()
	if err != nil {
		return g.setStatusError(StateMissing, g.downloader.ConfiguredSource(), err)
	}

	source := plan.Source
	g.setStatus(Status{
		State:  StateDownloading,
		DBPath: g.dbPath,
		Source: source,
	})

	if err := g.downloader.Download(ctx, plan); err != nil {
		return g.setStatusError(StateError, source, err)
	}

	if err := g.lookup.Load(); err != nil {
		return g.setStatusError(StateError, source, err)
	}

	g.setReadyStatus(source)
	return nil
}

func (g *Service) setReadyStatus(source Source) {
	g.setStatus(Status{
		State:     StateReady,
		DBPath:    g.dbPath,
		Source:    source,
		UpdatedAt: g.lookup.UpdatedAt(),
	})
}

func (g *Service) setStatus(status Status) {
	if status.UpdatedAt == nil {
		now := time.Now()
		status.UpdatedAt = &now
	}
	g.statusMu.Lock()
	g.status = status
	g.statusMu.Unlock()
}

func (g *Service) setStatusError(state State, source Source, err error) error {
	g.setStatus(Status{
		State:     state,
		DBPath:    g.dbPath,
		Source:    source,
		LastError: err.Error(),
	})
	return err
}
