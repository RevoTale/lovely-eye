package downloader

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/lovely-eye/server/internal/geoip"
)

type DownloadPlan struct {
	Source        geoip.Source
	CandidateURLs []string
}

type Downloader struct {
	dbPath string

	downloadURL       string
	maxMindLicenseKey string
	httpClient        *http.Client
}

func New(cfg geoip.Config) *Downloader {
	return &Downloader{
		dbPath:            cfg.DBPath,
		downloadURL:       strings.TrimSpace(cfg.DownloadURL),
		maxMindLicenseKey: strings.TrimSpace(cfg.MaxMindLicenseKey),
		httpClient:        &http.Client{Timeout: 30 * time.Second},
	}
}

func (d *Downloader) HasDownloadSource() bool {
	return d.ConfiguredSource() != geoip.SourceUnknown
}

func (d *Downloader) ConfiguredSource() geoip.Source {
	switch {
	case d.downloadURL != "":
		if isDBIPDownloadURL(d.downloadURL) {
			return geoip.SourceDBIP
		}
		return geoip.SourceDownloadURL
	case d.maxMindLicenseKey != "":
		return geoip.SourceMaxMind
	default:
		return geoip.SourceUnknown
	}
}

func (d *Downloader) Download(ctx context.Context, plan DownloadPlan) error {
	return d.downloadAndInstall(ctx, plan)
}
