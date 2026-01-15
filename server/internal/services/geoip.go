package services

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
)

const (
	geoIPStateDisabled    = "disabled"
	geoIPStateMissing     = "missing"
	geoIPStateDownloading = "downloading"
	geoIPStateReady       = "ready"
	geoIPStateError       = "error"
)

type GeoIPConfig struct {
	DBPath            string
	DownloadURL       string
	MaxMindLicenseKey string
}

type GeoIPStatus struct {
	State     string
	DBPath    string
	Source    string
	LastError string
	UpdatedAt *time.Time
}

// GeoIPService provides IP geolocation capabilities.
type GeoIPService struct {
	reader *geoip2.Reader
	mu     sync.RWMutex

	dbPath            string
	downloadURL       string
	maxMindLicenseKey string

	status   GeoIPStatus
	statusMu sync.RWMutex

	enabled    bool
	downloadMu sync.Mutex
	httpClient *http.Client
}

// NewGeoIPService creates a new GeoIP service.
// The database is downloaded on demand when country tracking is enabled.
func NewGeoIPService(cfg GeoIPConfig) (*GeoIPService, error) {
	service := &GeoIPService{
		dbPath:            cfg.DBPath,
		downloadURL:       cfg.DownloadURL,
		maxMindLicenseKey: cfg.MaxMindLicenseKey,
		httpClient:        &http.Client{Timeout: 30 * time.Second},
	}
	service.setStatus(GeoIPStatus{
		State:  geoIPStateDisabled,
		DBPath: cfg.DBPath,
	})
	return service, nil
}

func (g *GeoIPService) SetEnabled(enabled bool) {
	g.statusMu.Lock()
	g.enabled = enabled
	if !enabled {
		g.status = GeoIPStatus{
			State:  geoIPStateDisabled,
			DBPath: g.dbPath,
		}
	}
	g.statusMu.Unlock()
}

func (g *GeoIPService) Status() GeoIPStatus {
	g.statusMu.RLock()
	defer g.statusMu.RUnlock()
	return g.status
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

	if g.dbPath == "" {
		return g.setStatusError(geoIPStateMissing, "GeoIP database path is not configured")
	}

	if g.hasReader() {
		g.setStatus(GeoIPStatus{
			State:  geoIPStateReady,
			DBPath: g.dbPath,
			Source: "file",
		})
		return nil
	}

	if g.fileExists(g.dbPath) {
		if err := g.loadReader(); err == nil {
			g.setStatus(GeoIPStatus{
				State:  geoIPStateReady,
				DBPath: g.dbPath,
				Source: "file",
			})
			return nil
		}
	}

	if !g.downloadConfigured() {
		return g.setStatusError(geoIPStateMissing, "GeoIP download source is not configured")
	}

	g.downloadMu.Lock()
	defer g.downloadMu.Unlock()

	if g.hasReader() {
		g.setStatus(GeoIPStatus{
			State:  geoIPStateReady,
			DBPath: g.dbPath,
			Source: "file",
		})
		return nil
	}

	if err := g.downloadAndLoad(ctx); err != nil {
		return err
	}

	g.setStatus(GeoIPStatus{
		State:  geoIPStateReady,
		DBPath: g.dbPath,
		Source: g.downloadSource(),
	})
	return nil
}

func (g *GeoIPService) downloadAndLoad(ctx context.Context) error {
	downloadURLs, source, err := g.resolveDownloadURL()
	if err != nil {
		return g.setStatusError(geoIPStateMissing, err.Error())
	}

	g.setStatus(GeoIPStatus{
		State:  geoIPStateDownloading,
		DBPath: g.dbPath,
		Source: source,
	})

	var downloadErr error
	for _, url := range downloadURLs {
		downloadErr = g.downloadDatabase(ctx, url)
		if downloadErr == nil {
			downloadErr = nil
			break
		}
	}
	if downloadErr != nil {
		return g.setStatusError(geoIPStateError, downloadErr.Error())
	}

	if err := g.loadReader(); err != nil {
		return g.setStatusError(geoIPStateError, err.Error())
	}

	return nil
}

func (g *GeoIPService) downloadDatabase(ctx context.Context, url string) error {
	if err := os.MkdirAll(filepath.Dir(g.dbPath), 0o755); err != nil {
		return fmt.Errorf("create GeoIP directory: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build GeoIP request: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("download GeoIP database: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download GeoIP database: unexpected status %s", resp.Status)
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(g.dbPath), "geoip-*.download")
	if err != nil {
		return fmt.Errorf("create GeoIP temp file: %w", err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("save GeoIP download: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close GeoIP download: %w", err)
	}

	if g.isTarGz(url, resp.Header.Get("Content-Type")) {
		return g.extractTarGz(tmpFile.Name())
	}

	if g.isGzip(url, resp.Header.Get("Content-Type")) {
		return g.extractGzip(tmpFile.Name())
	}

	return g.replaceDatabase(tmpFile.Name())
}

func (g *GeoIPService) extractTarGz(archivePath string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("open GeoIP archive: %w", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("read GeoIP gzip: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read GeoIP archive: %w", err)
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		if !strings.HasSuffix(header.Name, ".mmdb") {
			continue
		}

		tmpPath := g.dbPath + ".tmp"
		out, err := os.Create(tmpPath)
		if err != nil {
			return fmt.Errorf("create GeoIP database: %w", err)
		}
		if _, err := io.Copy(out, tarReader); err != nil {
			_ = out.Close()
			return fmt.Errorf("write GeoIP database: %w", err)
		}
		if err := out.Close(); err != nil {
			return fmt.Errorf("close GeoIP database: %w", err)
		}
		return g.replaceDatabase(tmpPath)
	}

	return fmt.Errorf("GeoIP archive did not contain an .mmdb file")
}

func (g *GeoIPService) replaceDatabase(tempPath string) error {
	if err := os.Chmod(tempPath, 0o644); err != nil {
		return fmt.Errorf("set GeoIP database permissions: %w", err)
	}
	if err := os.Rename(tempPath, g.dbPath); err != nil {
		return fmt.Errorf("move GeoIP database: %w", err)
	}
	return nil
}

func (g *GeoIPService) isTarGz(url, contentType string) bool {
	if strings.HasSuffix(url, ".tar.gz") || strings.HasSuffix(url, ".tgz") {
		return true
	}
	return strings.Contains(contentType, "gzip")
}

func (g *GeoIPService) isGzip(url, contentType string) bool {
	if strings.HasSuffix(url, ".gz") {
		return true
	}
	return strings.Contains(contentType, "gzip")
}

func (g *GeoIPService) extractGzip(archivePath string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("open GeoIP archive: %w", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("read GeoIP gzip: %w", err)
	}
	defer gzipReader.Close()

	tmpPath := g.dbPath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create GeoIP database: %w", err)
	}
	if _, err := io.Copy(out, gzipReader); err != nil {
		_ = out.Close()
		return fmt.Errorf("write GeoIP database: %w", err)
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("close GeoIP database: %w", err)
	}
	return g.replaceDatabase(tmpPath)
}

func (g *GeoIPService) resolveDownloadURL() ([]string, string, error) {
	if g.downloadURL != "" {
		source := "download-url"
		if strings.Contains(g.downloadURL, "db-ip.com") {
			source = "dbip"
		}
		return g.expandDownloadURLs(g.downloadURL), source, nil
	}
	if g.maxMindLicenseKey == "" {
		return nil, "", errors.New("GeoIP download URL or MaxMind license key is required")
	}
	url := fmt.Sprintf("https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=%s&suffix=tar.gz", g.maxMindLicenseKey)
	return []string{url}, "maxmind", nil
}

func (g *GeoIPService) expandDownloadURLs(url string) []string {
	if !strings.Contains(url, "dbip-country-lite.mmdb.gz") {
		return []string{url}
	}

	now := time.Now()
	candidates := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		month := now.AddDate(0, -i, 0).Format("2006-01")
		candidate := strings.Replace(url, "dbip-country-lite.mmdb.gz", fmt.Sprintf("dbip-country-lite-%s.mmdb.gz", month), 1)
		candidates = append(candidates, candidate)
	}
	return candidates
}

func (g *GeoIPService) downloadConfigured() bool {
	return g.downloadURL != "" || g.maxMindLicenseKey != ""
}

func (g *GeoIPService) downloadSource() string {
	if g.downloadURL != "" {
		return "download-url"
	}
	if g.maxMindLicenseKey != "" {
		return "maxmind"
	}
	return ""
}

func (g *GeoIPService) loadReader() error {
	reader, err := geoip2.Open(g.dbPath)
	if err != nil {
		return fmt.Errorf("open GeoIP database: %w", err)
	}
	g.mu.Lock()
	if g.reader != nil {
		_ = g.reader.Close()
	}
	g.reader = reader
	g.mu.Unlock()
	return nil
}

func (g *GeoIPService) hasReader() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.reader != nil
}

func (g *GeoIPService) fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func (g *GeoIPService) setStatus(status GeoIPStatus) {
	now := time.Now()
	status.UpdatedAt = &now
	g.statusMu.Lock()
	g.status = status
	g.statusMu.Unlock()
}

func (g *GeoIPService) setStatusError(state, message string) error {
	g.setStatus(GeoIPStatus{
		State:     state,
		DBPath:    g.dbPath,
		Source:    g.downloadSource(),
		LastError: message,
	})
	return errors.New(message)
}

// GetCountry returns the ISO country code for an IP address.
func (g *GeoIPService) GetCountry(ipStr string) string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "Unknown"
	}

	// Skip private/local IPs
	if isPrivateIP(ip) {
		return "Local"
	}

	if g.reader == nil {
		return "Unknown"
	}

	record, err := g.reader.Country(ip)
	if err != nil {
		return "Unknown"
	}

	if record.Country.IsoCode == "" {
		return "Unknown"
	}

	return record.Country.IsoCode
}

// GetCountryName returns the country name for an IP address.
func (g *GeoIPService) GetCountryName(ipStr string) string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "Unknown"
	}

	if isPrivateIP(ip) {
		return "Local Network"
	}

	if g.reader == nil {
		return "Unknown"
	}

	record, err := g.reader.Country(ip)
	if err != nil {
		return "Unknown"
	}

	if record.Country.Names["en"] == "" {
		return "Unknown"
	}

	return record.Country.Names["en"]
}

// Close closes the GeoIP database reader.
func (g *GeoIPService) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.reader != nil {
		return g.reader.Close()
	}
	return nil
}

// isPrivateIP checks if an IP is private/local.
func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsPrivate() {
		return true
	}
	return false
}
