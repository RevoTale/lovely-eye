package services

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang/v2"
	"github.com/oschwald/maxminddb-golang"
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

type GeoIPCountry struct {
	Code string
	Name string
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

	countriesMu        sync.RWMutex
	countries          []GeoIPCountry
	countriesUpdatedAt *time.Time
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
		g.clearCountriesCache()
	}
	g.statusMu.Unlock()
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

	if status.UpdatedAt != nil {
		if cached := g.getCachedCountries(*status.UpdatedAt); cached != nil {
			return filterCountries(cached, search), nil
		}
	}

	reader, err := maxminddb.Open(g.dbPath)
	if err != nil {
		return nil, fmt.Errorf("open GeoIP database: %w", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			slog.Error("failed to close maxminddb reader", "error", err)
		}
	}()

	type countryRecord struct {
		Country struct {
			ISOCode string            `maxminddb:"iso_code"`
			Names   map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
		RegisteredCountry struct {
			ISOCode string            `maxminddb:"iso_code"`
			Names   map[string]string `maxminddb:"names"`
		} `maxminddb:"registered_country"`
	}

	network := reader.Networks(maxminddb.SkipAliasedNetworks)
	seen := make(map[string]string)
	var record countryRecord
	for network.Next() {
		_, err := network.Network(&record)
		if err != nil {
			return nil, fmt.Errorf("read GeoIP network: %w", err)
		}

		code := record.Country.ISOCode
		name := record.Country.Names["en"]
		if code == "" {
			code = record.RegisteredCountry.ISOCode
			name = record.RegisteredCountry.Names["en"]
		}
		if code == "" {
			continue
		}
		if name == "" {
			name = code
		}
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = name
	}
	if err := network.Err(); err != nil {
		return nil, fmt.Errorf("read GeoIP networks: %w", err)
	}

	countries := make([]GeoIPCountry, 0, len(seen))
	for code, name := range seen {
		countries = append(countries, GeoIPCountry{
			Code: code,
			Name: name,
		})
	}
	sort.Slice(countries, func(i, j int) bool {
		if countries[i].Name == countries[j].Name {
			return countries[i].Code < countries[j].Code
		}
		return countries[i].Name < countries[j].Name
	})

	if status.UpdatedAt != nil {
		g.setCachedCountries(countries, *status.UpdatedAt)
	}
	return filterCountries(countries, search), nil
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
	if err := os.MkdirAll(filepath.Dir(g.dbPath), 0o750); err != nil {
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

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
	file, err := os.Open(archivePath) // #nosec G304 -- archivePath is constructed from validated dbPath
	if err != nil {
		return fmt.Errorf("open GeoIP archive: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("failed to close archive file", "error", err)
		}
	}()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("read GeoIP gzip: %w", err)
	}
	defer func() {
		if err := gzipReader.Close(); err != nil {
			slog.Error("failed to close gzip reader", "error", err)
		}
	}()

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
		out, err := os.Create(tmpPath) // #nosec G304 -- tmpPath is constructed from validated dbPath
		if err != nil {
			return fmt.Errorf("create GeoIP database: %w", err)
		}
		// #nosec G110 -- GeoIP database files are from trusted sources (MaxMind)
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
	if err := os.Chmod(tempPath, 0o600); err != nil {
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
	file, err := os.Open(archivePath) // #nosec G304 -- archivePath is constructed from validated dbPath
	if err != nil {
		return fmt.Errorf("open GeoIP archive: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("failed to close archive file", "error", err)
		}
	}()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("read GeoIP gzip: %w", err)
	}
	defer func() {
		if err := gzipReader.Close(); err != nil {
			slog.Error("failed to close gzip reader", "error", err)
		}
	}()

	tmpPath := g.dbPath + ".tmp"
	out, err := os.Create(tmpPath) // #nosec G304 -- tmpPath is constructed from validated dbPath
	if err != nil {
		return fmt.Errorf("create GeoIP database: %w", err)
	}
	// #nosec G110 -- GeoIP database files are from trusted sources (MaxMind)
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

func (g *GeoIPService) getCachedCountries(updatedAt time.Time) []GeoIPCountry {
	g.countriesMu.RLock()
	defer g.countriesMu.RUnlock()
	if g.countriesUpdatedAt == nil || !g.countriesUpdatedAt.Equal(updatedAt) {
		return nil
	}
	result := make([]GeoIPCountry, len(g.countries))
	copy(result, g.countries)
	return result
}

func (g *GeoIPService) setCachedCountries(countries []GeoIPCountry, updatedAt time.Time) {
	g.countriesMu.Lock()
	g.countries = make([]GeoIPCountry, len(countries))
	copy(g.countries, countries)
	g.countriesUpdatedAt = &updatedAt
	g.countriesMu.Unlock()
}

func (g *GeoIPService) clearCountriesCache() {
	g.countriesMu.Lock()
	g.countries = nil
	g.countriesUpdatedAt = nil
	g.countriesMu.Unlock()
}

func filterCountries(countries []GeoIPCountry, search string) []GeoIPCountry {
	query := strings.TrimSpace(strings.ToLower(search))
	if query == "" {
		result := make([]GeoIPCountry, len(countries))
		copy(result, countries)
		return result
	}

	filtered := make([]GeoIPCountry, 0, len(countries))
	for _, country := range countries {
		code := strings.ToLower(country.Code)
		name := strings.ToLower(country.Name)
		if strings.Contains(code, query) || strings.Contains(name, query) {
			filtered = append(filtered, country)
		}
	}
	return filtered
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

var ErrNoDBReader = errors.New("no IP reader")
var UnknownCountry = Country{
		Name: "Unknown",
		ISOCode: "-",
}
var LocalNetworkCountry = Country{
	Name: "Local Network",
	ISOCode: "-",
}
type Country struct {
	Name string
	ISOCode string
}
func (g *GeoIPService) ResolveCountry(ipStr string) (Country,error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	ip,err := netip.ParseAddr(ipStr)
	if err != nil {
		return UnknownCountry, fmt.Errorf("failed parse country IP: %s",err.Error())
	}

	if ip.IsPrivate() {
		return LocalNetworkCountry, nil
	}

	if g.reader == nil {
		return UnknownCountry, ErrNoDBReader
	}

	record, err := g.reader.Country(ip)
	if err != nil {
		return UnknownCountry, fmt.Errorf("failed to get country: %s",err.Error())
	}

	return Country{
		Name: record.Country.Names.English,
		ISOCode: record.Country.ISOCode,
	},nil
}

func (g *GeoIPService) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.reader != nil {
		if err := g.reader.Close(); err != nil {
			return fmt.Errorf("failed to close geoip reader: %w", err)
		}
	}
	return nil
}

